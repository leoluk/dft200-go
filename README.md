# dft200-go

Quick and dirty Go CLI for the Sportstech DFT200 standing desk treadmill.

Requires a running bluetoothd and Go 1.11+.

Inspired by https://github.com/machinekoder/deskfit.

## Disclaimer

This is an unofficial hobby project with no warranties or liability, implied or otherwise.

Using a treadmill is dangerous. A few notes by a fellow DFT200 owner:

- The Bluetooth feature is unauthenticated. Anyone in range can control your treadmill.

- I haven't figured out the right command for emergency stop, and I found no references
  while reverse engineering the Android application. **Always keep the hardware remote control close to you**,
  especially at higher speeds - pressing the start button *twice* will trigger an emergency stop.
  
- Be careful when playing with the Bluetooth interface, you might brick your treadmill.
  There are commands for reading and (presumably) writing the EEPROM.
  
- Do not install the Android application - it requests tons of dangerous permissions for no
  good reason, and at least one variant of it is packed by a APK packer commonly used by malware.
  The DFT200 manual recommends to install the APK from their website,
  which is generally a bad idea™.

## Setup 

Install bluetoothd via your distro package manager and start it:

    systemctl enable bluetooth.service
    systemctl start bluetooth.service

Find MAC address using `bluetoothctl`:

    scan on
    devices
    connect <mac>
    
Alternatively, you can use a bluetoothd GUI like Gnome's native applet, blueberry (GTK) or bluedevil (KDE).

Install to `$GOPATH/bin`:

    go install github.com/leoluk/dft200-go/cmd/dft-cli

## Usage

Start treadmill (or continue, if paused):

    dft-cli -addr <mac> -start
    
Pause treadmill:

    dft-cli -addr <mac> -pause
    
Stop treadmill (you don't usually need this, 
it will stop after you leave it paused for a few minutes):

    dft-cli -addr <mac> -stop

Set speed (10-80 for levels 1-8):

    dft-cli -addr <mac> -speed <n>

## Example i3 config

```
set $tmMac <mac>

bindsym $mod+Ctrl+1 exec ~/go/bin/dft-cli -addr $tmMac -speed 10
bindsym $mod+Ctrl+2 exec ~/go/bin/dft-cli -addr $tmMac -speed 20
bindsym $mod+Ctrl+3 exec ~/go/bin/dft-cli -addr $tmMac -speed 30
bindsym $mod+Ctrl+4 exec ~/go/bin/dft-cli -addr $tmMac -speed 40
bindsym $mod+Ctrl+5 exec ~/go/bin/dft-cli -addr $tmMac -speed 50
bindsym $mod+Ctrl+6 exec ~/go/bin/dft-cli -addr $tmMac -speed 60
bindsym $mod+Ctrl+7 exec ~/go/bin/dft-cli -addr $tmMac -speed 70
bindsym $mod+Ctrl+8 exec ~/go/bin/dft-cli -addr $tmMac -speed 80
bindsym $mod+Ctrl+BackSpace exec ~/go/bin/dft-cli -addr $tmMac -pause
bindsym $mod+Ctrl+Return exec ~/go/bin/dft-cli -addr $tmMac -start
```
 
