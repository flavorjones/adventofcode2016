#! /usr/bin/env ruby

Node = Struct.new(:x, :y, :size, :used, keyword_init: true) do
  class << self
    def from_desc(desc)
      x, y, size, used = desc.scan(/\d+/).map(&:to_i)
      new(x: x, y: y, size: size, used: used)
    end
  end

  def initialize(**)
    super
  end

  def avail
    size - used
  end
end

class Grid
  attr_reader :width, :height

  def initialize(input)
    @nodes = []
    @width = 0
    @height = 0
    scan(input)
  end

  def at(x, y)
    @nodes[y]&.[](x)
  end

  def adjacent_to(node)
    [
      at(node.x + 1, node.y),
      at(node.x, node.y + 1),
    ].compact.tap do |nodes|
      nodes << at(node.x - 1, node.y) if node.x > 0
      nodes << at(node.x, node.y - 1) if node.y > 0
    end
  end

  def scan(input)
    input.each_line do |line|
      next unless line.start_with?("/dev/grid/node-")

      node = Node.from_desc(line)
      @nodes[node.y] ||= []
      @nodes[node.y][node.x] = node
      @width = [@width, node.x + 1].max
      @height = [@height, node.y + 1].max
    end
  end

  def viable_pairs
    @nodes.flatten.permutation(2).select do |a, b|
      a.used > 0 && a.used <= b.avail
    end
  end

  def inspect
    @nodes.map do |row|
      row.map do |node|
        if node.x == 0 && node.y == 0
          "x"
        elsif node.x == @width - 1 && node.y == 0
          "G"
        elsif node.used == 0
          "_"
        elsif node.used > 100
          "#"
        else
          "."
        end
      end.join
    end.join("\n")
  end
end

if ARGV[0] != "run"
  require "minitest/autorun"

  module Minitest::Assertions
    def assert_includes_exactly(expected, actual)
      assert_equal(expected.size, actual.size)
      expected.each { |e| assert_includes(actual, e) }
    end
  end

  describe Node do
    it "has attributes" do
      node = Node.new(x: 123, y: 345, size: 456, used: 135)
      assert_equal(123, node.x)
      assert_equal(345, node.y)
      assert_equal(456, node.size)
      assert_equal(135, node.used)
      assert_equal(456-135, node.avail)
    end

    it "can be constructed from a string description" do
      node = Node.from_desc("/dev/grid/node-x2-y20    85T   69T    16T   81%")
      assert_equal(2, node.x)
      assert_equal(20, node.y)
      assert_equal(85, node.size)
      assert_equal(69, node.used)
      assert_equal(16, node.avail)
    end
  end

  describe Grid do
    let(:input) { <<~INPUT }
      root@ebhq-gridcenter# df -h
      Filesystem              Size  Used  Avail  Use%
      /dev/grid/node-x0-y0     92T   68T    24T   73%
      /dev/grid/node-x0-y1     87T   73T    14T   83%
      /dev/grid/node-x0-y2     89T   64T    25T   71%
      /dev/grid/node-x0-y3     91T   64T    27T   70%
      /dev/grid/node-x1-y0     88T   65T    23T   73%
      /dev/grid/node-x1-y1     94T   69T    25T   73%
      /dev/grid/node-x1-y2     85T   70T    15T   82%
      /dev/grid/node-x1-y3    507T  493T    14T   97%
      /dev/grid/node-x2-y0     92T   71T    21T   77%
      /dev/grid/node-x2-y1     91T   69T    22T   75%
      /dev/grid/node-x2-y2     90T   70T    20T   77%
      /dev/grid/node-x2-y3    502T  490T    12T   97%
    INPUT

    it "can be constructed from a string description" do
      grid = Grid.new(input)
      assert_equal(3, grid.width)
      assert_equal(4, grid.height)

      node = grid.at(1, 2)
      assert_equal(85, node.size)
    end

    it "can find adjacent nodes" do
      grid = Grid.new(input)

      expected = [grid.at(1, 0), grid.at(0, 1), grid.at(2, 1), grid.at(1, 2)]
      assert_includes_exactly(expected, grid.adjacent_to(grid.at(1, 1)))

      expected = [grid.at(1, 2), grid.at(0, 3), grid.at(2, 3)]
      assert_includes_exactly(expected, grid.adjacent_to(grid.at(1, 3)))

      expected = [grid.at(0, 1), grid.at(1, 0)]
      assert_includes_exactly(expected, grid.adjacent_to(grid.at(0, 0)))
    end

    it "can find viable pairs" do
      grid = Grid.new(input)
      assert_empty(grid.viable_pairs)

      grid.at(1, 1).used = 30
      assert_includes_exactly(
        [[grid.at(0, 2), grid.at(1, 1)], [grid.at(0, 3), grid.at(1, 1)]],
        grid.viable_pairs,
      )
    end
  end
else
  input = File.read(File.join(__dir__, "day22.txt"))
  grid = Grid.new(input)
  puts "grid: #{grid.width}x#{grid.height}"

  # part 1
  puts "part 1: #{grid.viable_pairs.size}"

  # part 2
  puts grid.inspect
end
