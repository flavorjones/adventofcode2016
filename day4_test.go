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
)

type SortableNameComponent struct {
	element     byte
	occurrences uint
}

type SortableNameComponents []SortableNameComponent

func (c SortableNameComponents) Len() int      { return len(c) }
func (c SortableNameComponents) Swap(j, k int) { c[j], c[k] = c[k], c[j] }
func (c SortableNameComponents) Less(j, k int) bool {
	if c[j].occurrences == c[k].occurrences {
		return c[j].element < c[k].element
	}
	return c[j].occurrences > c[k].occurrences
}

type Room struct {
	descriptor string
}

var roomSectorIdRe = regexp.MustCompile(`-(\d+)\[`)
var roomNameRe = regexp.MustCompile(`([-a-z]+)-\d`)
var roomNameIgnore = "-"[0]
var roomDescribedChecksumRe = regexp.MustCompile(`\[(.*)\]`)
var blankStringRe = regexp.MustCompile(`^\s*$`)

func NewRoom(descriptor string) *Room {
	return &Room{descriptor}
}

func (r Room) sectorID() int {
	match := roomSectorIdRe.FindStringSubmatch(r.descriptor)
	sectorId, _ := strconv.ParseInt(match[1], 10, 16)
	return int(sectorId)
}

func (r Room) name() string {
	return roomNameRe.FindStringSubmatch(r.descriptor)[1]
}

func (r Room) describedChecksum() string {
	return roomDescribedChecksumRe.FindStringSubmatch(r.descriptor)[1]
}

func (r Room) valid() bool {
	described := r.describedChecksum()
	actual := r.nameChecksum()
	return described == actual
}

func (r Room) nameChecksum() string {
	// build a map
	byteCount := make(map[byte]uint)
	for _, char := range []byte(r.name()) {
		if char != "-"[0] {
			if current, ok := byteCount[char]; ok {
				byteCount[char] = current + 1
			} else {
				byteCount[char] = 1
			}
		}
	}

	// convert to an array of SortableNameComponent
	components := make(SortableNameComponents, 0, len(byteCount))
	for element, occurrences := range byteCount {
		components = append(components, SortableNameComponent{element, occurrences})
	}
	sort.Sort(components)

	// assemble the string
	rval := make([]byte, 0, len(components))
	for _, component := range components {
		rval = append(rval, component.element)
	}
	return string(rval[0:5])
}

var _ = Describe("Day4", func() {
	Describe("Room", func() {
		room1 := NewRoom("aaaaa-bbb-z-y-x-123[abxyz]")
		room2 := NewRoom("a-b-c-d-e-f-g-h-987[abcde]")
		room3 := NewRoom("not-a-real-room-404[oarel]")
		room4 := NewRoom("totally-real-room-200[decoy]")

		Describe("#valid", func() {
			It("can detect decoys", func() {
				Expect(room1.valid()).To(BeTrue())
				Expect(room2.valid()).To(BeTrue())
				Expect(room3.valid()).To(BeTrue())
				Expect(room4.valid()).To(BeFalse())
			})
		})

		Describe("#sectorID", func() {
			It("returns an integer sector ID from the descriptor", func() {
				Expect(room1.sectorID()).To(Equal(123))
				Expect(room2.sectorID()).To(Equal(987))
				Expect(room3.sectorID()).To(Equal(404))
				Expect(room4.sectorID()).To(Equal(200))
			})
		})

		Describe("#name", func() {
			It("returns the encrypted name from the descriptor", func() {
				Expect(room1.name()).To(Equal("aaaaa-bbb-z-y-x"))
				Expect(room2.name()).To(Equal("a-b-c-d-e-f-g-h"))
				Expect(room3.name()).To(Equal("not-a-real-room"))
				Expect(room4.name()).To(Equal("totally-real-room"))
			})
		})

		Describe("#describedChecksum", func() {
			It("returns the checksum from the descriptor", func() {
				Expect(room1.describedChecksum()).To(Equal("abxyz"))
				Expect(room2.describedChecksum()).To(Equal("abcde"))
				Expect(room3.describedChecksum()).To(Equal("oarel"))
				Expect(room4.describedChecksum()).To(Equal("decoy"))
			})
		})
	})

	Describe("the puzzle", func() {
		It("star 1", func() {
			data, _ := ioutil.ReadFile("day4.txt")
			sum := 0
			for _, line := range strings.Split(string(data), "\n") {
				if blankStringRe.MatchString(line) {
					continue
				}
				if room := NewRoom(line); room.valid() {
					sum += room.sectorID()
				}
			}
			fmt.Println("sum is ", sum)
		})
	})
})
