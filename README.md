# Firestack

Firestack is a userspace TCP/UDP connection monitor, firewall, DNS resolver, and multi-hop [WireGuard](https://github.com/wireguard/wireguard-go) client for Android.

It is a fork of [celzero/firestack](https://github.com/celzero/firestack), which powers Rethink DNS. [gVisor/netstack](https://github.com/google/gvisor/tree/go/pkg/tcpip) provides a SOCKS-like interface (similar to [badvpn's tun2socks](https://github.com/ambrop72/badvpn)) for TCP/UDP over a TUN device. The project started as a hard fork of Google's [outline-go-tun2socks](https://github.com/Jigsaw-Code/outline-go-tun2socks).

This fork is used by [DNS Control](https://github.com/hidefbi/dns-control).

## What's different here

- The DNS listener has rank hints: `OnPrequery` and `DomainOpts`, with `DNSSummary` carrying the rank.
- The flow listener gets a connection id (`cid`) on `Preflow`, `Flow` and `Inflow`. It is echoed back in `Postflow`.
- Only the Android build is kept up to date. Other targets in the Makefile are inherited and not tested here.

## DNS

Firestack supports DNS over HTTPS, DNS over TLS, Oblivious DNS over HTTPS, DNS over WireGuard / SOCKS5 / Tor, DNSCrypt, and plain old DNS upstreams.

## WireGuard

Firestack runs WireGuard in userspace. When running *multiple* WireGuard tunnels at once, only ICMP, DNS, TCP and UDP are forwarded through them. ARP / IGMP / SCTP / RTP and other IP protocols are *not* forwarded to WireGuard tunnels.

Firestack supports multi-hop / multi-relay WireGuard, where multiple tunnels can be chained together, provided that the outer tunnel (hop/relay) can route to the inner tunnel's (exit) endpoint.

WireGuard integration was sponsored by [FOSS United](https://fossunited.org/grants); and Multi-hop / Multi-relay WireGuard by [FLOSS/fund](https://floss.fund/).

## Use it

Firestack is published through [JitPack](https://jitpack.io/#hidefbi/firestack). Builds are per commit, so pin a commit hash (26727b7 is the one DNS Control uses).

```gradle
// project build.gradle
allprojects {
    repositories {
        ...
        maven { url 'https://jitpack.io' }
        ...
    }
}

// app build.gradle
dependencies {
    ...
    implementation 'com.github.hidefbi:firestack:26727b7@aar'
    ...
}
```

Build logs are at `https://jitpack.io/com/github/hidefbi/firestack/<commit>/build.log`.

This fork does not publish to Maven Central.

## API

The API isn't stable and isn't documented. The easiest reference is the integration in [DNS Control](https://github.com/hidefbi/dns-control): see `app/src/main/java/com/celzero/bravedns/service/BraveVPNService.kt`.

## Build

Only Android is supported.

### Prerequisites

- macOS or Linux
- make
- Go 1.27.3 (see `go.mod`)
- JDK 17, for the Gradle side
- Android SDK, platform 36
- Android NDK 28.2.13676358 for local builds

JitPack builds with NDK `30.0.16248370` and platform 37 (see `jitpack.yml`).

### Android

```bash
# sdkmanager is at ~/Android/Sdk/cmdline-tools/latest/bin/sdkmanager
~/Android/Sdk/cmdline-tools/latest/bin/sdkmanager "platforms;android-36" "ndk;28.2.13676358"
export ANDROID_HOME=~/Android/Sdk ANDROID_NDK_HOME=~/Android/Sdk/ndk/28.2.13676358
```

[gomobile](https://pkg.go.dev/golang.org/x/mobile/cmd/gobind) is installed by `make` as needed.

## Make

```
# creates build/intra/tun2socks.aar and tun2socks-sources.jar
make clean && make intra
```

The first build needs network. The modules are vendored, but `gomobile bind` still reads the module cache, so the Makefile sets `GOFLAGS=-mod=mod`.

For a debug AAR with symbols, use `make intradebug`.

To pull the native libraries out of the AAR:
```bash
unzip build/intra/tun2socks.aar 'jni/*' -d build/android
```

## License

Mozilla Public License 2.0. See [LICENSE](LICENSE).
