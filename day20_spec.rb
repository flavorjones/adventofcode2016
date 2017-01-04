require 'rspec'

class FirewallRules
  class Blacklist
    RULE_RE = /(\d+)-(\d+)/

    attr_reader :ranges

    def initialize rules
      @ranges = rules.each_line.map do |rule|
        match = RULE_RE.match rule
        match[1].to_i..match[2].to_i
      end
    end

    def in number
      ranges.
        select {|range| range.cover? number}.
        sort_by {|r| -r.end}
    end
  end

  attr_reader :blacklist, :upper_limit

  def initialize blacklist, upper_limit=4294967295
    @upper_limit = upper_limit
    @blacklist = Blacklist.new blacklist
  end

  def lowest_unblocked_addy
    j = 0
    loop do
      high_range = blacklist.in(j).first
      return j if high_range.nil?
      j = high_range.end + 1
    end
  end

  def number_of_allowed_addresses
    j = 0
    count = 0
    while j <= upper_limit
      high_range = blacklist.in(j).first
      if high_range.nil?
        count += 1
        j += 1
      else
        j = high_range.end + 1
      end
    end
    count
  end
end


describe FirewallRules do
  describe FirewallRules::Blacklist do
    let(:blacklist) { <<~EOB }
      5-8
      0-2
      4-7
      6-12
    EOB

    let(:b) { FirewallRules::Blacklist.new blacklist }

    describe ".new" do
      it "creates a Range to match each rule" do
        expect(b.ranges).to match_array([5..8, 0..2, 4..7, 6..12])
      end
    end

    describe "#in" do
      it "returns all the ranges in which the number exists" do
        expect(b.in(2)).to match_array([0..2])
        expect(b.in(5)).to match_array([4..7, 5..8])
        expect(b.in(8)).to match_array([5..8, 6..12])
      end

      it "returns those ranges sorted inversely by the range.end" do
        expect(b.in(6)).to eq([6..12, 5..8, 4..7])
      end

      context "number is not in any range" do
        it "returns [] " do
          expect(b.in(3)).to eq([])
        end
      end
    end
  end

  context "when initialized with blacklists" do
    let(:blacklist) { <<~EOB }
      5-8
      0-2
      4-7
    EOB

    let(:fr) { FirewallRules.new blacklist, 9 }

    it "finds the lowest unblocked IP" do
      expect(fr.lowest_unblocked_addy).to eq(3)
    end

    it "counts the number of allowed IPs" do
      expect(fr.number_of_allowed_addresses).to eq(2)
    end
  end
end

describe "the puzzle" do
  let(:blacklist) { File.read("day20.txt") }
  let(:fr) { FirewallRules.new blacklist }

  it "star 1" do
    addy = fr.lowest_unblocked_addy
    puts "day 20 star 1: lowest IP is #{addy}"
  end

  it "star 2" do
    n = fr.number_of_allowed_addresses
    puts "day 20 star 2: total allowed addies is #{n}"
  end
end
