# coding: utf-8
require 'rspec'

class Room
  attr :descriptor

  def initialize descriptor
    @descriptor = descriptor
  end

  def sectorID
    match = /-(\d+)\[/.match descriptor
    match[1].to_i
  end

  def name
    match = /([-a-z]+)-\d/.match descriptor
    match[1]
  end

  def describedChecksum
    match = /\[(.*)\]/.match descriptor
    match[1]
  end

  def valid?
    describedChecksum == nameChecksum
  end

  def nameChecksum
    char_count = {}
    char_count.default = 0
    name.each_char do |char|
      next if char == "-"
      char_count[char] += 1
    end
    char_count.sort_by do |key, value|
      [-value, key]
    end.map(&:first).take(5).join
  end
end


describe Room do
  let(:room1) { Room.new("aaaaa-bbb-z-y-x-123[abxyz]") }
	let(:room2) { Room.new("a-b-c-d-e-f-g-h-987[abcde]") }
	let(:room3) { Room.new("not-a-real-room-404[oarel]") }
	let(:room4) { Room.new("totally-real-room-200[decoy]") }

  describe "#valid" do
    it "can detect decoys" do
      expect(room1.valid?).to be true
      expect(room2.valid?).to be true
      expect(room3.valid?).to be true
      expect(room4.valid?).to be false
    end
  end

	describe "#sectorID" do
		it "returns an integer sector ID from the descriptor" do
			expect(room1.sectorID).to eq(123)
			expect(room2.sectorID).to eq(987)
			expect(room3.sectorID).to eq(404)
			expect(room4.sectorID).to eq(200)
    end
  end

	describe "#name" do
		it "returns the encrypted name from the descriptor" do
			expect(room1.name).to eq("aaaaa-bbb-z-y-x")
			expect(room2.name).to eq("a-b-c-d-e-f-g-h")
			expect(room3.name).to eq("not-a-real-room")
			expect(room4.name).to eq("totally-real-room")
    end
  end

	describe "#describedChecksum" do
		it "returns the checksum from the descriptor" do
			expect(room1.describedChecksum).to eq("abxyz")
			expect(room2.describedChecksum).to eq("abcde")
			expect(room3.describedChecksum).to eq("oarel")
			expect(room4.describedChecksum).to eq("decoy")
    end
  end
end
