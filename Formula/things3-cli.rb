class Things3Cli < Formula
  desc "CLI for Things 3"
  homepage "https://github.com/ossianhempel/things3-cli"
  version "0.3.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/ossianhempel/things3-cli/releases/download/v0.3.0/things-0.3.0-darwin-arm64.tar.gz"
      sha256 "a6464403832457a935d7b7131d358802d5b751f4c7d25a06725b7bd3dad4e68f"
    else
      url "https://github.com/ossianhempel/things3-cli/releases/download/v0.3.0/things-0.3.0-darwin-amd64.tar.gz"
      sha256 "235dc9f0a60499d392457114a52b83d5643045ee1020d9f3fa0a97444c6a4197"
    end
  end

  def install
    bin.install "things"
  end

  test do
    system "#{bin}/things", "--version"
  end
end
