# Maintainer: vinegm <info AT vinegm D0T com>
pkgname='ditui-bin'
pkgver=1.0.0
pkgrel=1
pkgdesc="Morse code TUI trainer"
arch=('x86_64' 'aarch64')
url="https://github.com/vinegm/ditui"
license=('MIT')
provides=('ditui')
conflicts=('ditui')
options=('!strip' '!debug')
source_x86_64=("ditui-${pkgver}-x86_64.tar.gz::${url}/releases/download/v${pkgver}/ditui_${pkgver}_linux_amd64.tar.gz")
source_aarch64=("ditui-${pkgver}-aarch64.tar.gz::${url}/releases/download/v${pkgver}/ditui_${pkgver}_linux_arm64.tar.gz")
sha256sums_x86_64=('3c3041f7bbf4ca59be644caebdff3ba40578b945c07f7bcb22b8b2c988e22ab7')
sha256sums_aarch64=('af9268000468b97d5c5be715dc34bd9bac2971a6b7635f57674bd7be8828e92f')

package() {
    install -Dm755 ditui "$pkgdir/usr/bin/ditui"
    install -Dm644 LICENSE "$pkgdir/usr/share/licenses/$pkgname/LICENSE"
}
