require 'rspec'

class KineticSculpture
  class Disc
    SPEC_RE =
      /Disc #(\d+) has (\d+) positions; at time=(\d+), it is at position (\d+)/

    attr :positions, :initial_offset

    def initialize spec
      matches = SPEC_RE.match spec
      raise "bad time #{matches[3]}" unless matches[3] == "0"

      @positions = matches[2].to_i
      @initial_offset = matches[4].to_i
    end

    def in_position_at? t
      (t + initial_offset) % positions == 0
    end
  end

  attr :discs

  def initialize disc_arrangement
    @discs = disc_arrangement.split("\n").map do |spec|
      Disc.new spec
    end
  end

  def ball_drop_time
    t0 = 0
    loop do
      t = t0
      break if discs.all? do |disc|
        t += 1
        disc.in_position_at? t
      end

      t0 += 1
    end
    return t0
  end
end


describe KineticSculpture do
  let(:arrangement) { <<~EOA }
    Disc #1 has 5 positions; at time=0, it is at position 4.
    Disc #2 has 2 positions; at time=0, it is at position 1.
  EOA

  describe KineticSculpture::Disc do
    Disc = KineticSculpture::Disc

    describe ".new" do
      it "parses its spec" do
        disc = Disc.new "Disc #99 has 2 positions; at time=0, it is at position 1."
        expect(disc.positions).to eq(2)
        expect(disc.initial_offset).to eq(1)
      end
    end

    describe "#in_position_at?" do
      it "returns true if the disc will be in position at time t" do
        disc = Disc.new "Disc #99 has 4 positions; at time=0, it is at position 1."
        expect(disc.in_position_at?(0)).to be(false)
        expect(disc.in_position_at?(1)).to be(false)
        expect(disc.in_position_at?(2)).to be(false)
        expect(disc.in_position_at?(3)).to be(true)
        expect(disc.in_position_at?(4)).to be(false)
      end
    end
  end

  describe ".new" do
    it "creates a disc for each line in the arrangement" do
      ks = KineticSculpture.new arrangement
      expect(ks.discs.length).to eq(arrangement.split("\n").length)
    end
  end

  describe "#ball_drop_time" do
    it "finds the first time at which the ball can be dropped" do
      ks = KineticSculpture.new arrangement
      expect(ks.ball_drop_time).to eq(5)
    end
  end
end

describe "the puzzle" do
  let(:disc_arrangement) { <<~EOA }
    Disc #1 has 17 positions; at time=0, it is at position 1.
    Disc #2 has 7 positions; at time=0, it is at position 0.
    Disc #3 has 19 positions; at time=0, it is at position 2.
    Disc #4 has 5 positions; at time=0, it is at position 0.
    Disc #5 has 3 positions; at time=0, it is at position 0.
    Disc #6 has 13 positions; at time=0, it is at position 5.
  EOA

  it "star 1" do
    ks = KineticSculpture.new disc_arrangement
    puts "star 1: drop ball at #{ks.ball_drop_time}"
  end

  it "star 2" do
    da = disc_arrangement + "Disc #7 has 11 positions; at time=0, it is at position 0."
    ks = KineticSculpture.new da
    puts "star 2: drop ball at #{ks.ball_drop_time}"
  end
end
