# GOPATH and ANDROID_HOME needs to be set
# gomobile & gobind needs to be installed in $GOPATH/bin

mkdir -p build/android
gomobile bind -target=android -ldflags=-extldflags=-Wl,-soname,libgojni.so -tags="android experimental signrpc walletrpc chainrpc invoicesrpc routerrpc backuprpc peerrpc submarineswaprpc breezbackuprpc" -o build/android/breez.aar github.com/breez/breez/bindings
