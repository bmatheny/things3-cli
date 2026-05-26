class Things3Cli < Formula
  desc "CLI for Things 3"
  homepage "https://github.com/ossianhempel/things3-cli"
  version "0.3.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/ossianhempel/things3-cli/releases/download/v0.3.0/things-0.3.0-darwin-arm64.tar.gz"
      sha256 "4a87d9544f3a1357c1b7a335a3432bd2865a7f3181ae211260e1484ec3dd48b0"
    else
      url "https://github.com/ossianhempel/things3-cli/releases/download/v0.3.0/things-0.3.0-darwin-amd64.tar.gz"
      sha256 "92d074d9d8de81e314b217a3a0f9de507f9ccede182725657bd30a229107e983"
    end
  end

  def install
    bin.install "things"
  end

  test do
    system "#{bin}/things", "--version"
  end
end
