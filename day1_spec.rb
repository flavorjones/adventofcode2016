require 'rspec'

class Position
  HEADINGS = [
    NORTH = [0, 1],
    EAST  = [1, 0],
    SOUTH = [0, -1],
    WEST  = [-1, 0],
  ]

  attr :x, :y, :heading

  def initialize
    @x = 0
    @y = 0
    @heading = NORTH
  end

  def move step
    match = /([LR])([0-9]+)/.match(step)
    turn match[1]
    walk match[2].to_i
    self
  end

  def taxicab_geometry
    x.abs + y.abs
  end

  def position
    [x, y]
  end

  private

  def turn direction
    @heading = case direction
               when "L"
                 case heading
                 when NORTH then WEST
                 when WEST then SOUTH
                 when SOUTH then EAST
                 else NORTH
                 end
               else
                 case heading
                 when NORTH then EAST
                 when EAST then SOUTH
                 when SOUTH then WEST
                 else NORTH
                 end
               end
  end

  def walk distance
    @x += heading[0] * distance
    @y += heading[1] * distance
  end
end

describe Position do
  context "heading" do
    it { expect(Position.new.move("R0").heading).to eq(Position::EAST) }
    it { expect(Position.new.move("R0").move("R0").heading).to eq(Position::SOUTH) }
    it { expect(Position.new.move("R0").move("R0").move("R0").heading).to eq(Position::WEST) }
    it { expect(Position.new.move("R0").move("R0").move("R0").move("R0").heading).to eq(Position::NORTH) }
    it { expect(Position.new.move("L0").heading).to eq(Position::WEST) }
    it { expect(Position.new.move("L0").move("L0").heading).to eq(Position::SOUTH) }
    it { expect(Position.new.move("L0").move("L0").move("L0").heading).to eq(Position::EAST) }
    it { expect(Position.new.move("L0").move("L0").move("L0").move("L0").heading).to eq(Position::NORTH) }
  end

  context "distance" do
    it { expect(Position.new.move("R2").position).to eq([2, 0]) }
    it { expect(Position.new.move("R0").move("R2").position).to eq([0, -2]) }
    it { expect(Position.new.move("R0").move("R0").move("R2").position).to eq([-2, 0]) }
    it { expect(Position.new.move("R0").move("R0").move("R0").move("R2").position).to eq([0, 2]) }
  end
end

class GridPath
  attr :pathstring

  def initialize pathstring
    @pathstring = pathstring
  end

  def distance
    pathstring.split(/, ?/).inject(Position.new) do |position, step|
      position.move step
    end.taxicab_geometry
  end
end

describe GridPath do
  it { expect(GridPath.new("R2, L3").distance).to eq 5 }
  it { expect(GridPath.new("R2, R2, R2").distance).to eq 2 }
  it { expect(GridPath.new("R5, L5, R5, R3").distance).to eq 12 }
end

path = "L1, L5, R1, R3, L4, L5, R5, R1, L2, L2, L3, R4, L2, R3, R1, L2, R5, R3, L4, R4, L3, R3, R3, L2, R1, L3, R2, L1, R4, L2, R4, L4, R5, L3, R1, R1, L1, L3, L2, R1, R3, R2, L1, R4, L4, R2, L189, L4, R5, R3, L1, R47, R4, R1, R3, L3, L3, L2, R70, L1, R4, R185, R5, L4, L5, R4, L1, L4, R5, L3, R2, R3, L5, L3, R5, L1, R5, L4, R1, R2, L2, L5, L2, R4, L3, R5, R1, L5, L4, L3, R4, L3, L4, L1, L5, L5, R5, L5, L2, L1, L2, L4, L1, L2, R3, R1, R1, L2, L5, R2, L3, L5, L4, L2, L1, L2, R3, L1, L4, R3, R3, L2, R5, L1, L3, L3, L3, L5, R5, R1, R2, L3, L2, R4, R1, R1, R3, R4, R3, L3, R3, L5, R2, L2, R4, R5, L4, L3, L1, L5, L1, R1, R2, L1, R3, R4, R5, R2, R3, L2, L1, L5"

puts GridPath.new(path).distance
