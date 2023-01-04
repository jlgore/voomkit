# VoomKit

## Divoom Pixoo64 integration with HomeKit

This module is built upon the shoulders of [Brutella's HAP](https://github.com/brutella/hap), [Xxandev's homekit](https://github.com/xxandev/homekit), and [Jurjevic's pixoo64](https://github.com/jurjevic/pixoo64/). It's goal is to expose the Pixoo64 from Divoom as a HomeKit compatible TV.

If you want to play around with this, install `docker` and `make`. From there you can `make build`, `make run`, or `make all` (to build and run at once.)

Currently Working:
- [x] Turn Display off/on.
- [x] Change television inputs to different channels on the Pixoo64. (switching from inputs in the upper ranges causes the channel to not change.)

Todo:
- [ ] Set custom Divoom channels via CLI flag or OS environment variables.
- [ ] Write messages on Divoom on events.
