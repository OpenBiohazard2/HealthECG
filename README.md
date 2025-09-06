# HealthECG
Render health ECG (electrocardiogram) animation

<img src="https://github.com/OpenBiohazard2/HealthECG/blob/main/screenshots/screenshot.png" alt="screenshot" width="400" height="300" />

## Usage

```bash
go run .
```

The program will open a window showing 5 different ECG states with animated waveforms.

## About

This program shows how the health ECG animation was implemented in the original Resident Evil 2 code. This animation always renders when the player opens up the inventory screen in the top part of the screen.

The program displays two views for each ECG state:
- **Main View (left)**: Animated scrolling ECG with gradient fading
- **Overview (right)**: Complete static ECG curve with viewport indicator

The animation is rendered by drawing a contiguous subset of the ECG on each frame. Each subset of the ECG forms a curve, but the curve is composed of multiple vertical lines instead of one single curved line. The position (yOffset from the top left corner) and height of the vertical lines is hardcoded in an array, which depends on the player's status, e.g. fine, caution, danger, etc. The color and gradient is also defined for each player status in RGB form, where the minimum value is 0 and the maximum value is 255. For example, when the player's status is fine, the color will be green. The gradient is used to blend in the ECG line with a black background.

When the player opens up the inventory, the ECG starts drawing the first vertical line at position x = 0. The code will try to draw multiple vertical lines to the left, where each successive line to the left is darker than the previous line. The color of the left vertical lines is defined by subtracting the gradient from the adjacent line on the right. For example, if the current line color is (20, 255, 20) and the gradient is (1, 8, 1), the first line on the left will have color (19, 247, 19), the second line will have color (18, 239, 18), the third line will have color (17, 231, 17) and so on. The code can draw up to 32 vertical lines at most to the left of the initial vertical line. If any vertical line goes to the left side of the ECG at x = 0, the line won't be rendered.

For each frame, the current vertical line position will be incremented by 1, starting from x = 0 and going all the way up to x = 127. This increment makes the animation smooth, since you gradually see the ECG line move to the right per frame.
