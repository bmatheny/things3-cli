class Things3Cli < Formula
  desc "CLI for Things 3"
  homepage "https://github.com/bmatheny/things3-cli"
  version "0.4.0-beta"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/bmatheny/things3-cli/releases/download/v0.4.0-beta/things-0.4.0-beta-darwin-arm64.tar.gz"
      sha256 "912931b0264224e70d668d91e77d35eb24da30a92bcbd6a3baba8c17aa2d6220"
    else
      url "https://github.com/bmatheny/things3-cli/releases/download/v0.4.0-beta/things-0.4.0-beta-darwin-amd64.tar.gz"
      sha256 "c710e965afe662f4c793555b29e98c3d86c8c17d530405ac863801e28c140dc2"
    end
  end

  def install
    bin.install "things"
  end

  test do
    system "#{bin}/things", "--version"
  end
end
