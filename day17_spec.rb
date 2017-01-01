require 'rspec'

class VaultFinder
  FLOOR_WIDTH = 4
  FLOOR_HEIGHT = 4

  DIRECTIONS = %w[U D L R]
  OPEN_VALUES = %w[b c d e f]

  Position = Struct.new(:x, :y) do
    def directions
      DIRECTIONS.dup.tap do |directions|
        directions.delete "L" if x == 0
        directions.delete "U" if y == 0
        directions.delete "R" if x == FLOOR_WIDTH-1
        directions.delete "D" if y == FLOOR_HEIGHT-1
      end
    end

    def new_position direction
      # skipping wall detection
      case direction
      when "U" then Position.new(x, y-1)
      when "D" then Position.new(x, y+1)
      when "L" then Position.new(x-1, y)
      when "R" then Position.new(x+1, y)
      end
    end
  end

  VAULT_POSITION = Position.new(3, 3)
  START_POSITION = Position.new(0, 0)

  attr :passcode

  def initialize passcode
    @passcode = passcode
  end

  def shortest_path
    paths = {"" => START_POSITION}
    loop do
      break if paths.empty?
      new_paths = {}
      paths.each do |passcode_path, position|
        (position.directions & directions(passcode_path)).each do |direction|
          new_path = passcode_path + direction
          new_position = position.new_position(direction)
          return new_path if new_position == VAULT_POSITION
          new_paths[new_path] = new_position
        end
      end
      paths = new_paths
    end
  end

  def longest_path
    longest_path_length = 0
    longest_path = nil
    paths = {"" => START_POSITION}
    loop do
      break if paths.empty?
      new_paths = {}
      paths.each do |passcode_path, position|
        (position.directions & directions(passcode_path)).each do |direction|
          new_path = passcode_path + direction
          new_position = position.new_position(direction)
          if new_position == VAULT_POSITION
            if longest_path_length < new_path.length
              longest_path_length = new_path.length
              longest_path = new_path
            end
            next
          end
          new_paths[new_path] = new_position
        end
      end
      paths = new_paths
    end
    longest_path
  end

  def directions path
    Digest::MD5.hexdigest(passcode + path).chars[0..3].zip(DIRECTIONS).map do |data|
      key, direction = *data
      direction if OPEN_VALUES.include? key
    end.compact
  end
end


describe VaultFinder do
  describe VaultFinder::Position do
    Position = VaultFinder::Position
    describe "#directions" do
      it "returns an array of available directions for a room" do
        expect(Position.new(0, 0).directions).to match_array(%w[D R])
        expect(Position.new(0, 1).directions).to match_array(%w[U D R])
        expect(Position.new(1, 0).directions).to match_array(%w[R D L])
        expect(Position.new(1, 1).directions).to match_array(%w[U D L R])
        expect(Position.new(0, 3).directions).to match_array(%w[R U])
        expect(Position.new(3, 0).directions).to match_array(%w[L D])
      end
    end

    describe "#new_position" do
      it "returns the Position for the position in that direction" do
        expect(Position.new(1, 1).new_position("R")).to eq(Position.new(2, 1))
        expect(Position.new(1, 1).new_position("D")).to eq(Position.new(1, 2))
        expect(Position.new(1, 1).new_position("L")).to eq(Position.new(0, 1))
        expect(Position.new(1, 1).new_position("U")).to eq(Position.new(1, 0))
      end
    end
  end

  describe ".directions" do
    it "returns an array of available directions for a passcode" do
      expect(VaultFinder.new("hijkl").directions("")).to match_array(%w[U D L])
      expect(VaultFinder.new("hijkl").directions("D")).to match_array(%w[U L R])
    end
  end

  describe "#shortest_path" do
    it "detects when no path is possible" do
      vf = VaultFinder.new "hijkl"
      expect(vf.shortest_path).to be_nil
    end

    it "detects the shortest path" do
      expect(VaultFinder.new("ihgpwlah").shortest_path).to eq("DDRRRD")
      expect(VaultFinder.new("kglvqrro").shortest_path).to eq("DDUDRLRRUDRD")
      expect(VaultFinder.new("ulqzkmiv").shortest_path).to eq("DRURDRUDDLLDLUURRDULRLDUUDDDRR")
    end
  end

  describe "#longest_path" do
    it "detects when no path is possible" do
      vf = VaultFinder.new "hijkl"
      expect(vf.longest_path).to be_nil
    end

    it "detects the longest path" do
      expect(VaultFinder.new("ihgpwlah").longest_path.length).to eq(370)
      expect(VaultFinder.new("kglvqrro").longest_path.length).to eq(492)
      expect(VaultFinder.new("ulqzkmiv").longest_path.length).to eq(830)
    end
  end
end

describe "the puzzle" do
  it "star 1" do
    path = VaultFinder.new("dmypynyp").shortest_path
    puts "star 1: shortest path is #{path}"
  end

  it "star 2" do
    path = VaultFinder.new("dmypynyp").longest_path
    puts "star 2: longest path is #{path}"
    puts "star 2: longest path length is #{path.length}"
  end
end
