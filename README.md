# vision-summary
models that summarize information from underlying vision models 

## Example Config

### for count-classifier
```
{
  "count_thresholds": {
    "high": 1000,
    "none": 0,
    "low": 10,
    "medium": 20
  },
  "detector_name": "vision-1",
  "chosen_labels": {
    "person": 0.3
  }
}
```

### for count-sensor
```
{
  "count_thresholds": {
    "high": 1000,
    "none": 0,
    "low": 10,
    "medium": 20
  },
  "detector_name": "vision-1",
  "camera_name": "camera-1",
  "poll_frequency_hz": 0.5,
  "chosen_labels": {
    "person": 0.3
  }
}
```
