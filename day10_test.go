package adventofcode2016_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type BotDistributionRule struct {
	low     int
	lowBot  bool
	high    int
	highBot bool
}

type Bot struct {
	id        int
	chips     []int
	input     chan int
	powerDown chan bool
	inputMap  []chan int
	rule      BotDistributionRule
}

func NewBot(id int, rule BotDistributionRule) Bot {
	return Bot{
		id,
		make([]int, 0, 2),
		make(chan int),
		make(chan bool),
		nil,
		rule,
	}
}

func (b *Bot) PowerUp() {
	go func() {
		fmt.Println("Bot", b.id, "is powered up")
		for {
			select {
			case <-b.powerDown:
				fmt.Println("Bot", b.id, "is powering down")
				return
			case input := <-b.input:
				if len(b.chips) >= 2 {
					fmt.Println("Bot", b.id, "IGNORES input", input)
				} else {
					fmt.Println("Bot", b.id, "accepts input", input)
					b.chips = append(b.chips, input)
				}
			default:
				time.Sleep(10 * time.Millisecond)
			}

			if len(b.chips) == 2 && len(b.inputMap) > 0 {
				sort.Ints(b.chips)
				fmt.Println("Bot", b.id, "handling", b.chips)
				if b.rule.lowBot {
					b.inputMap[b.rule.low] <- b.chips[0]
				} else {
					fmt.Println("Bot", b.id, "delivered", b.chips[0], "to output", b.rule.low)
				}
				if b.rule.highBot {
					b.inputMap[b.rule.high] <- b.chips[1]
				} else {
					fmt.Println("Bot", b.id, "delivered", b.chips[1], "to output", b.rule.high)
				}
				b.chips = b.chips[:0]
			}
		}
	}()
}

func (b Bot) PowerDown() {
	b.powerDown <- true
}

type BotMaster struct {
	rules []string
	bots  []Bot
}

func NewBotMaster(rules []string) BotMaster {
	return BotMaster{rules, make([]Bot, 2)}
}

func createDistributionRule(match []string) BotDistributionRule {
	var lowBot, highBot bool
	if match[2] == "bot" {
		lowBot = true
	} else {
		lowBot = false
	}
	lowIndex, _ := strconv.Atoi(match[3])
	if match[4] == "bot" {
		highBot = true
	} else {
		highBot = false
	}
	highIndex, _ := strconv.Atoi(match[5])
	return BotDistributionRule{lowIndex, lowBot, highIndex, highBot}
}

var botRuleRe = regexp.MustCompile(`bot (\d+) gives low to (\w+) (\d+) and high to (\w+) (\d+)`)

func (bm *BotMaster) CreateBots() {
	for _, rule := range bm.rules {
		if match := botRuleRe.FindStringSubmatch(rule); match != nil {
			botNum, _ := strconv.Atoi(match[1])
			distributionRule := createDistributionRule(match)

			if len(bm.bots) <= botNum {
				bigger := make([]Bot, botNum+1)
				copy(bigger, bm.bots)
				bm.bots = bigger
			}
			fmt.Println("creating bot", botNum, len(bm.bots))
			bm.bots[botNum] = NewBot(botNum, distributionRule)
		}
	}

	for jbot, _ := range bm.bots {
		if bm.bots[jbot].input == nil {
			bm.bots = bm.bots[:jbot]
			break
		}
	}

	inputMap := make([]chan int, len(bm.bots))
	for jbot, bot := range bm.bots {
		inputMap[jbot] = bot.input
	}

	for jbot, _ := range bm.bots {
		bm.bots[jbot].inputMap = inputMap
	}
}

func (bm BotMaster) SeedBot() {
	re := regexp.MustCompile(`value (\d+) goes to bot (\d+)`)
	for _, rule := range bm.rules {
		if match := re.FindStringSubmatch(rule); match != nil {
			value, _ := strconv.Atoi(match[1])
			botNum, _ := strconv.Atoi(match[2])
			bm.bots[botNum].input <- value
		}
	}
}

func (bm *BotMaster) StartBots() {
	bm.CreateBots()

	for jbot, _ := range bm.bots {
		bm.bots[jbot].PowerUp()
		defer bm.bots[jbot].PowerDown()
	}

	bm.SeedBot()

	time.Sleep(time.Second * 2)
}

var _ = Describe("Day10", func() {
	Describe("Bot", func() {
		Context("when given an inputMap and a rule", func() {
			It("delivers chips to the right place", func() {
				bot0 := NewBot(0, BotDistributionRule{1, true, 0, false})
				bot1 := NewBot(1, BotDistributionRule{})

				inputMap := []chan int{bot0.input, bot1.input}
				bot0.inputMap = inputMap

				bot0.PowerUp()
				defer bot0.PowerDown()
				bot1.PowerUp()
				defer bot1.PowerDown()

				bot0.input <- 10
				bot0.input <- 20

				Eventually(func() int {
					return len(bot0.chips)
				}).Should(Equal(0))

				Eventually(func() []int {
					return bot1.chips
				}).Should(Equal([]int{10}))
			})
		})
	})

	Describe("BotMaster", func() {

		rawData := `value 5 goes to bot 2
bot 2 gives low to bot 1 and high to bot 0
value 3 goes to bot 1
bot 1 gives low to output 1 and high to bot 0
bot 0 gives low to output 2 and high to output 0
value 2 goes to bot 2`
		rules := strings.Split(string(rawData), "\n")

		Describe("#CreateFleet", func() {
			It("creates a fleet based on rules", func() {
				bm := NewBotMaster(rules)

				bm.StartBots()

				Expect(len(bm.bots)).To(Equal(3))
				Expect(len(bm.bots[0].inputMap)).To(Equal(3))
			})
		})
	})

	Describe("the puzzle", func() {
		rawData, _ := ioutil.ReadFile("day10.txt")
		rules := strings.Split(string(rawData), "\n")

		It("star 1", func() {
			bm := NewBotMaster(rules)
			bm.StartBots()
		})
	})
})
