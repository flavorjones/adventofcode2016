require 'rspec'
require 'ruby-prof'
require 'set'

def Profile &block
  if ENV["PROFILE"]
    result = RubyProf.profile(&block)
    printer = RubyProf::CallStackPrinter.new(result)
    printer.print(File.open("day11_callstack.html", "w"))
  else
    block.call
  end
end


class RTFA # RadioisotopeTestingFacility Artifact
  include Comparable

  MICROCHIP = :microchip
  GENERATOR = :generator

  attr :element, :artifact_type, :hashable

  def initialize element, artifact_type
    @element = element.to_sym
    @artifact_type = artifact_type.to_sym
    @hashable = [element, artifact_type]
  end

  def <=> other
    @hashable <=> other.hashable
  end

  def to_s
    sprintf "%.2s-%.1s", element, artifact_type
  end

  def inspect
    to_s
  end

  def match? other
    return false if artifact_type == other.artifact_type
    return other.element == element
  end

  def hash
    @hash ||= @hashable.hash
  end
end

def RTFG *args
  RTFG.new(*args)
end

module RTFG # RadioisotopeTestingFacility Generator
  def self.new element
    RTFA.new element, RTFA::GENERATOR
  end
end

def RTFM *args
  RTFM.new(*args)
end

module RTFM # RadioisotopeTestingFacility Microchip
  def self.new element
    RTFA.new element, RTFA::MICROCHIP
  end
end

class RTF # RadioisotopeTestingFacility
  include Comparable

	MICROCHIP_RE = /\b(\w+)-compatible microchip/
	GENERATOR_RE = /\b(\w+) generator/
    
  def self.new_from_description description
    floors = []
    description.split("\n").each_with_index do |line, j|
      floors[j] = []
      line.scan(MICROCHIP_RE).each do |match|
        floors[j] << RTFM.new(match[0])
      end
      line.scan(GENERATOR_RE).each do |match|
        floors[j] << RTFG.new(match[0])
      end
    end
    new floors: floors
  end

  attr :floors, :elevator_pos, :hashable

  def initialize floors: [[],[],[],[]], elevator_pos: 0, sorted: false
    @floors = floors
    floors.map!(&:sort) unless sorted
    @elevator_pos = elevator_pos
    @hashable = [elevator_pos, floors]
  end

  def to_s
    output = []
    j = floors.length
    floors.reverse.each do |floor|
      elevator_eh = (elevator_pos == j-1)
      output << "F#{j}: #{elevator_eh ? "*" : " "} #{floor.join(', ')}"
      j -= 1
    end
    output << "\n"
    output.join("\n")
  end

  def valid?
    floors.each do |floor|
      floor.each do |artifact|
        next if artifact.artifact_type == RTFA::GENERATOR
        next if floor.any? { |other| artifact.match?(other) }
        return false if floor.any? { |other| other.artifact_type == RTFA::GENERATOR }
      end
    end
    true
  end

  def permutations
    elevator_combos = floors[elevator_pos].combination(1).to_a +
                      floors[elevator_pos].combination(2).to_a
    [].tap do |permutations|
      elevator_combos.each do |artifacts|
        if elevator_pos > 0
          new_floors = floors.map(&:dup)
          new_elevator_pos = elevator_pos - 1
          new_floors[elevator_pos] -= artifacts
          new_floors[elevator_pos].sort!
          new_floors[new_elevator_pos] += artifacts
          new_floors[new_elevator_pos].sort!
          permutations << RTF.new(floors: new_floors, elevator_pos: new_elevator_pos, sorted: true)
        end

        if elevator_pos < floors.length-1
          new_floors = floors.map(&:dup)
          new_elevator_pos = elevator_pos + 1
          new_floors[elevator_pos] -= artifacts
          new_floors[elevator_pos].sort!
          new_floors[new_elevator_pos] += artifacts
          new_floors[new_elevator_pos].sort!
          permutations << RTF.new(floors: new_floors, elevator_pos: new_elevator_pos, sorted: true)
        end
      end
    end
  end

  def valid_permutations
    permutations.select(&:valid?)
  end

  def <=> other
    @hashable <=> other.hashable
  end

  def eql? other
    hash == other.hash
  end

  def hash
    @hash ||= @hashable.hash # cached
  end

  def done?
    floors[-1].length > 0 && floors[0...-1].all?(&:empty?)
  end
end

class RTFTrip
  def initialize start
    @start = start
  end

  def solution_length
    all_permutations = Set.new
    all_permutations.add @start
    current_generation = [@start]
    jgen = 0

    loop do
      next_generation = current_generation.map do |rtf|
        rtf.valid_permutations.reject { |p| all_permutations.include?(p) }
      end.flatten.uniq

      puts "MIKE: gen #{jgen} has #{next_generation.length} new permutations"
      break if current_generation.any?(&:done?)

      next_generation.each { |p| all_permutations.add p }
      current_generation = next_generation
      jgen += 1
    end
    jgen
  end
end


describe RTF do
  let :description do
    <<~EOD
      The first floor contains a hydrogen-compatible microchip and a lithium-compatible microchip.
      The second floor contains a hydrogen generator.
      The third floor contains a lithium generator.
      The fourth floor contains nothing relevant.
    EOD
  end

  describe RTFA do
    describe "#match?" do
      context "for a microchip" do
        context "compared to a microchip" do
          it { expect(RTFM.new("a").match?(RTFM.new("a"))).to be_falsey }
          it { expect(RTFM.new("a").match?(RTFM.new("b"))).to be_falsey }
        end

        context "compared to a generator" do
          it { expect(RTFM.new("a").match?(RTFG.new("a"))).to be_truthy }
          it { expect(RTFM.new("a").match?(RTFG.new("b"))).to be_falsey }
        end
      end

      context "for a generator" do
        context "compared to a generator" do
          it { expect(RTFG.new("a").match?(RTFG.new("a"))).to be_falsey }
          it { expect(RTFG.new("a").match?(RTFG.new("b"))).to be_falsey }
        end

        context "compared to a microchip" do
          it { expect(RTFG.new("a").match?(RTFM.new("a"))).to be_truthy }
          it { expect(RTFG.new("a").match?(RTFM.new("b"))).to be_falsey }
        end
      end
    end
  end

  describe ".new_from_description" do
    it "takes a string description" do
      rtf = RTF.new_from_description description
      expect(rtf.floors[0]).to eq([RTFM.new("hydrogen"), RTFM.new("lithium")])
      expect(rtf.floors[1]).to eq([RTFG.new("hydrogen")])
      expect(rtf.floors[2]).to eq([RTFG.new("lithium")])
      expect(rtf.floors[3]).to eq([])
      expect(rtf.elevator_pos).to eq(0)
    end
  end

  describe ".new" do
    context "defaults" do
      it "has elevator position of zero" do
        expect(RTF.new.elevator_pos).to eq(0)
      end

      it "has empty floors" do
        expect(RTF.new.floors).to eq([[],[],[],[]])
      end
    end

    it "takes a floors array-of-arrays" do
      rtf = RTF.new floors: [[RTFM.new("a"), RTFM.new("b")], [RTFG.new("a")], [], []]
      expect(rtf.floors[0]).to eq([RTFM.new("a"), RTFM.new("b")])
      expect(rtf.floors[1]).to eq([RTFG.new("a")])
      expect(rtf.floors[2]).to eq([])
      expect(rtf.floors[3]).to eq([])
    end

    it "canonically sorts the artifacts on the floors" do
      rtf = RTF.new(floors: [
                      [RTFM.new("a"), RTFM.new("b"), RTFG.new("a"), RTFG.new("z")],
                      [RTFG.new("z"), RTFM.new("a"), RTFM.new("b"), RTFG.new("a")],
                    ])
      expect(rtf.floors).to(eq([
        [RTFG.new("a"), RTFM.new("a"), RTFM.new("b"), RTFG.new("z")],
        [RTFG.new("a"), RTFM.new("a"), RTFM.new("b"), RTFG.new("z")],
      ]))
    end

    it "takes an elevator position" do
      expect(RTF.new(elevator_pos: 3).elevator_pos).to eq(3)
    end
  end

  describe "#to_s" do
    it "emits something human-readable" do
      rtf = RTF.new_from_description description
      expect(rtf.to_s).to eq(<<~EOS)
        F4:   
        F3:   li-g
        F2:   hy-g
        F1: * hy-m, li-m

      EOS
    end
  end

  describe "#valid" do
    it "disallows a mismatched generator on a floor with a microchip" do
      expect(RTF.new(floors: [[RTFG.new("a"), RTFM.new("b")]]).valid?).to be_falsey
    end

    it "allows a matched generator on a floor with a microchip" do
      expect(RTF.new(floors: [[RTFG.new("a"), RTFM.new("a")]]).valid?).to be_truthy
    end

    it "allows a mismatched generator on a floor with a microchip/generator match" do
      expect(RTF.new(floors: [[RTFG.new("a"), RTFM.new("a"), RTFG.new("b")]]).valid?).to be_truthy
    end
  end

	describe "#permutations" do
		it "returns all possible next-step states" do
			initial = RTF.new(floors: [
                          [],
                          [RTFM.new("a"), RTFM.new("b"), RTFM.new("c")],
                          [],
                        ], elevator_pos: 1)

      expected = [
        RTF.new(floors: [[RTFM.new("a")], [RTFM.new("b"), RTFM.new("c")], []], elevator_pos: 0),
        RTF.new(floors: [[], [RTFM.new("b"), RTFM.new("c")], [RTFM.new("a")]], elevator_pos: 2),

        RTF.new(floors: [[RTFM.new("a"), RTFM.new("b")], [RTFM.new("c")], []], elevator_pos: 0),
        RTF.new(floors: [[], [RTFM.new("c")], [RTFM.new("a"), RTFM.new("b")]], elevator_pos: 2),

        RTF.new(floors: [[RTFM.new("a"), RTFM.new("c")], [RTFM.new("b")], []], elevator_pos: 0),
        RTF.new(floors: [[], [RTFM.new("b")], [RTFM.new("a"), RTFM.new("c")]], elevator_pos: 2),

        RTF.new(floors: [[RTFM.new("b")], [RTFM.new("a"), RTFM.new("c")], []], elevator_pos: 0),
        RTF.new(floors: [[], [RTFM.new("a"), RTFM.new("c")], [RTFM.new("b")]], elevator_pos: 2),

        RTF.new(floors: [[RTFM.new("b"), RTFM.new("c")], [RTFM.new("a")], []], elevator_pos: 0),
        RTF.new(floors: [[], [RTFM.new("a")], [RTFM.new("b"), RTFM.new("c")]], elevator_pos: 2),

        RTF.new(floors: [[RTFM.new("c")], [RTFM.new("a"), RTFM.new("b")], []], elevator_pos: 0),
        RTF.new(floors: [[], [RTFM.new("a"), RTFM.new("b")], [RTFM.new("c")]], elevator_pos: 2),
      ]

      actual = initial.permutations

      expect(actual).to match_array(expected)
    end
  end

  describe "#valid_permutations" do
    it "returns all valid next-step states" do
      initial = RTF.new floors: [[], [RTFM("a"), RTFG("a"), RTFM("c")], []], elevator_pos: 1

      expected = [
        RTF.new(floors: [[RTFM("a"), RTFG("a")], [RTFM("c")], []], elevator_pos: 0),
        RTF.new(floors: [[], [RTFM("c")], [RTFM("a"), RTFG("a")]], elevator_pos: 2),

        RTF.new(floors: [[RTFG("a")], [RTFM("a"), RTFM("c")], []], elevator_pos: 0),
        RTF.new(floors: [[], [RTFM("a"), RTFM("c")], [RTFG("a")]], elevator_pos: 2),

        RTF.new(floors: [[RTFM("c")], [RTFG("a"), RTFM("a")], []], elevator_pos: 0),
        RTF.new(floors: [[], [RTFG("a"), RTFM("a")], [RTFM("c")]], elevator_pos: 2),

        RTF.new(floors: [[RTFM("a"), RTFM("c")], [RTFG("a")], []], elevator_pos: 0),
        RTF.new(floors: [[], [RTFG("a")], [RTFM("a"), RTFM("c")]], elevator_pos: 2),
      ]

      actual = initial.valid_permutations

      expect(actual).to match_array(expected)
    end
  end

  describe "hash" do
    it "returns the same thing for similarly-structure RTFs" do
      rtf1 = RTF.new(floors: [[RTFM("c"), RTFM("a")], [], [RTFG("a")], [RTFG("c")]], elevator_pos: 2)
      rtf2 = RTF.new(floors: [[RTFM("c"), RTFM("a")], [], [RTFG("a")], [RTFG("c")]], elevator_pos: 2)
      expect(rtf1.hash).to eq(rtf2.hash)
    end
  end

  describe "eql?" do
    it "returns true for similary-structured RTFs" do
      rtf1 = RTF.new(floors: [[RTFM("c"), RTFM("a")], [], [RTFG("a")], [RTFG("c")]], elevator_pos: 2)
      rtf2 = RTF.new(floors: [[RTFM("c"), RTFM("a")], [], [RTFG("a")], [RTFG("c")]], elevator_pos: 2)
      expect(rtf1.eql?(rtf2)).to be(true)
    end
  end

  describe "#done?" do
    context "all artifacts are on the top floor" do
      it "returns true" do
        expect(RTF.new(floors: [[], [], [], [RTFM("a"), RTFG("a")]]).done?).to be_truthy
      end
    end

    context "some artifacts are not on the top floor" do
      it "returns false" do
        expect(RTF.new(floors: [[], [], [RTFM("b")], [RTFM("a"), RTFG("a")]]).done?).to be_falsey
        expect(RTF.new(floors: [[], [RTFM("b")], [], [RTFM("a"), RTFG("a")]]).done?).to be_falsey
        expect(RTF.new(floors: [[RTFM("b")], [], [], [RTFM("a"), RTFG("a")]]).done?).to be_falsey
      end
    end
  end

  describe "the described test" do
    it "finds a solution" do
      Profile do
        rtf = RTF.new_from_description description
        expect(RTFTrip.new(rtf).solution_length).to eq(11)
      end
    end
  end
end

describe "the puzzle" do
  describe "star 1" do
    let :description do
      <<~EOD
        The first floor contains a polonium generator, a thulium generator, a thulium-compatible microchip, a promethium generator, a ruthenium generator, a ruthenium-compatible microchip, a cobalt generator, and a cobalt-compatible microchip.
      	The second floor contains a polonium-compatible microchip and a promethium-compatible microchip.
      	The third floor contains nothing relevant.
      	The fourth floor contains nothing relevant.
      EOD
    end

    it do
      Profile do
        rtf = RTF.new_from_description description
        solution_length = RTFTrip.new(rtf).solution_length
        puts "star 1 solution length #{solution_length}"
      end
    end
  end

  describe "star 2" do
    let :description do
      <<~EOD
        The first floor contains a polonium generator, a thulium generator, a thulium-compatible microchip, a promethium generator, a ruthenium generator, a ruthenium-compatible microchip, a cobalt generator, a cobalt-compatible microchip, an elerium generator, an elerium-compatible microchip, a dilithium generator, and a dilithium-compatible microchip.
      	The second floor contains a polonium-compatible microchip and a promethium-compatible microchip.
      	The third floor contains nothing relevant.
      	The fourth floor contains nothing relevant.
      EOD
    end

    it do
      rtf = RTF.new_from_description description
      solution_length = RTFTrip.new(rtf).solution_length
      puts "star 1 solution length #{solution_length}"
    end
  end
end
