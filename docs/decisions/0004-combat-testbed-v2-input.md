# 0004 — Combat Testbed v2 Input Handling

**Date:** 2026-03-01
**Status:** Accepted

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>

## Decision

Rewrite combat testbed input handling with a two-mode design:

- **KR (Key Release) mode**: uses `held` map (set on press, cleared on release) and `pressed` map (one-shot accumulator using `IsRepeat` filtering). Zero allocations per tick.
- **FB (Fallback) mode**: timeout-based key tracking with edge detection via `prevHeld` snapshot. Timeout reduced from 150ms to 100ms.

A `[KR]`/`[FB]` diagnostic indicator in the HUD shows which mode is active.

## Context

After the bubbletea v2 migration, movement in the combat testbed felt "too far" — the character overshot when releasing a direction key.

### Root cause

bubbletea enables `ENABLE_VIRTUAL_TERMINAL_INPUT` on Windows. The ultraviolet terminal reader (`terminal_reader_windows.go`) then takes the VT input path, which only writes `kevent.Char` for `kevent.KeyDown` events — KEY_UP events are silently dropped.

Key releases (`KeyReleaseMsg`) only work if the terminal supports the Kitty keyboard protocol, which encodes release/repeat info as VT escape sequences that flow through the same KEY_DOWN path.

Without Kitty support, `SupportsEventTypes()` returns false, the testbed falls into fallback mode, and the 150ms timeout causes ~30px of coast (19% of the 160px arena width) per key release.

### Timeout tradeoff

| Timeout | Coast after release | Stutter on hold start |
|---------|--------------------|-----------------------|
| 150ms   | ~30px (5 ticks)    | ~100ms gap            |
| 100ms   | ~20px (3 ticks)    | ~150ms gap            |
| 66ms    | ~13px (2 ticks)    | ~184ms gap            |

**Choice: 100ms.** Cuts coast from 30px to 20px. The ~150ms stutter gap (during the ~250ms Windows keyboard repeat delay) is tolerable for a testbed and preferable to 19% overshoot.

## Key changes

- `pressed` map accumulates initial presses (`!msg.IsRepeat`) between ticks; cleared after `buildInput` reads it. Replaces `prevHeld` for one-shot detection in KR mode.
- `prevHeld` retained only for fallback mode edge detection.
- Key name constants (`keyLeft`, `keyRight`, `keyJump`, `keyAtk`, `keyReset`) eliminate magic strings.
- KR branch in `buildInput()` reads maps directly with no copies; uses `clear()` builtin.

## Alternatives considered

- **Only use fallback mode**: simpler but 30px coast is too much for a platformer.
- **Read Windows Console API directly**: would give native key releases but requires bypassing bubbletea's input layer — too invasive for a testbed.
- **Lower timeout to 66ms**: less coast but more noticeable stutter during keyboard repeat delay.
