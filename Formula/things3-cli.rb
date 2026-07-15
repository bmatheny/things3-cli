class Things3Cli < Formula
  desc "CLI for Things 3"
  homepage "https://github.com/ossianhempel/things3-cli"
  version "0.4.2-beta"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/ossianhempel/things3-cli/releases/download/v0.4.2-beta/things-0.4.2-beta-darwin-arm64.tar.gz"
      sha256 "620a78b6ba8248633ee828c662298891abe071849bb1f86c366c5ca84c0bc566"
    else
      url "https://github.com/ossianhempel/things3-cli/releases/download/v0.4.2-beta/things-0.4.2-beta-darwin-amd64.tar.gz"
      sha256 "5ae09226404d85bded3e9ed5936a5ca6db2ed13eae317405b587db297cb6ab65"
    end
  end

  def install
    bin.install "things"
  end

  test do
    system "#{bin}/things", "--version"
  end
end
