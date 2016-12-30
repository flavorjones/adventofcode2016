require 'rspec'
require 'set'

Position = Struct.new(:x, :y)

def Position x, y
  Position.new x, y
end


class Maze
  def initialize fav_number
    @fav_number = fav_number
  end

  def snapshot width, height
    output = []
    (0..height-1).each do |y|
      (0..width-1).each do |x|
        output << (wall?(Position(x, y)) ? "#" : ".")
      end
      output << "\n"
    end
    output.join
  end

  def wall? position
    x = position.x
    y = position.y
    
    value = x*x + 3*x + 2*x*y + y + y*y
    value += @fav_number

    binary = value.to_s(2)

    binary.chars.inject(0) do |count, char|
      char == "1" ? count+1 : count
    end.odd?
  end
end

class MazeSolver
  attr :maze, :start, :finish

  def initialize maze, start, finish
    @maze = maze
    @start = start
    @finish = finish
  end

  def solution_length
    visited_positions = Set.new
    visited_positions.add start
    current_step = [start]
    steps = 0

    loop do
      steps += 1

      next_step = current_step.map do |step|
        adjacent_spaces(step).reject { |p| visited_positions.include? p }
      end.flatten.uniq

      # puts "MIKE: gen #{steps} has #{next_step.length} new permutations"
      break if next_step.any? { |p| p == finish }

      next_step.each { |p| visited_positions.add p }
      current_step = next_step
    end
    steps
  end

  def adjacent_spaces position
    [].tap do |spaces|
      Position(position.x + 1, position.y).tap do |p|
        spaces << p unless maze.wall? p
      end      

      Position(position.x, position.y + 1).tap do |p|
        spaces << p unless maze.wall? p
      end      

      if position.x > 0
        Position(position.x - 1, position.y).tap do |p|
          spaces << p unless maze.wall? p
        end      
      end

      if position.y > 0
        Position(position.x, position.y - 1).tap do |p|
          spaces << p unless maze.wall? p
        end      
      end
    end
  end
end


describe Maze do
  describe "#wall?" do
    it "has walls in the right place" do
      maze = Maze.new(10)    
      expect(maze.wall?(Position(0, 0))).to be(false)
      expect(maze.wall?(Position(1, 0))).to be(true)
      expect(maze.wall?(Position(2, 0))).to be(false)
      expect(maze.wall?(Position(0, 1))).to be(false)
      expect(maze.wall?(Position(1, 1))).to be(false)
      expect(maze.wall?(Position(2, 1))).to be(true)
    end
  end

  describe "#snapshot" do
    it "has walls in the right place" do
      maze = Maze.new(10)
      expected = <<~EOT
        .#.####.##
        ..#..#...#
        #....##...
        ###.#.###.
        .##..#..#.
        ..##....#.
        #...##.###
      EOT
      actual = maze.snapshot(10, 7)
      expect(actual).to eq(expected)
    end
  end
end

describe MazeSolver do
  it "finds the shortest path" do
    solver = MazeSolver.new Maze.new(10), Position(1, 1), Position(7, 4)
    expect(solver.solution_length).to eq(11)
  end
end

describe "the puzzle" do
  describe "star 1" do
    it do
      maze = Maze.new 1364
      solver = MazeSolver.new maze, Position(1, 1), Position(31, 39)
      length = solver.solution_length
      puts "star 1: solution takes #{length} steps"
    end
  end
end
