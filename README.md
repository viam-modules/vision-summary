# vision-summary

This is a [Viam module](https://docs.viam.com/extend/modular-resources/) containing models that summarize information from underlying vision models.

## Getting started

First, [create a machine](https://docs.viam.com/manage/fleet/robots/#add-a-new-robot) in Viam.

Then, [add a module from the Viam Registry](https://docs.viam.com/modular-resources/configure/#add-a-module-from-the-viam-registry) and select the `viam:vision-summary:count-classifier` or `viam:vision-summary:count-sensor` model from the [`vision-summary` module](https://app.viam.com/module/viam/vision-summary).

## Configuration

### viam:vision-summary:count-classifier

To configure the `count-classifier` model, use the following template:

```
{
  "detector_name": <string>,
  "chosen_labels": {
    <label1>: <float>,
    <label2>: <float>
  }
}
```

#### Attributes

| Name | Type | Required? | Default | Description |
| ---- | ---- | --------- | --------| ------------ |
| `detector_name` | string | **Required** | | Name of the vision service to use as input. Must output a classifier tensor. |
| `chosen_labels` | string | **Required** | | Map of label names and required confidence values (between 0 and 1) to count the label in the summary. |
| `count_thresholds` | object | **Required** | | Maps summary category names to an integer value representing the (inclusive) upper bound of the range spanned by the category. Supports any number of custom category names. Each upper bound must be a unique integer. |

To configure `count_thresholds`, define at least one category. For instance, the following example defines categories spanning the following ranges:

- `none`: 0 to 10
- `low`: 11 to 20
- `medium`: 21 to 50
- `high`: 51 to 100

```json
"count_thresholds": {
  "none": 10,
  "low": 20,
  "medium": 50,
  "high": 100
}

### viam:vision-summary:count-sensor

To configure the `count-sensor` model, use the following template:

```
{
  "camera_name": <string>,
  "detector_name": <string>,
  "chosen_labels": {
    <label1>: <float>,
    <label2>: <float>
  },
  "count_thresholds": {
    "high": <int>,
    "none": <int>,
    "low": <int>,
    "medium": <int>
  }
}
```

#### Attributes

| Name | Type | Required? | Default | Description |
| ---- | ---- | --------- | --------| ------------ |
| `camera_name` | string | **Required** | | Camera name to use for video input. |
| `detector_name` | string | **Required** | | Name of the vision service to use as input. Must output a classifier tensor. |
| `chosen_labels` | string | **Required** | | Map of label names and required confidence values (between 0 and 1) to count the label in the summary. |
| `poll_frequency_hz` | object | Optional | | How many times to summarize per minute. |
| `count_thresholds` | object | **Required** | | Maps summary category names to an integer value representing the (inclusive) upper bound of the range spanned by the category. Supports any number of custom category names. Each upper bound must be a unique integer. |

To configure `count_thresholds`, define at least one category. For instance, the following example defines categories spanning the following ranges:

- `none`: 0 to 10
- `low`: 11 to 20
- `medium`: 21 to 50
- `high`: 51 to 100

```json
"count_thresholds": {
  "none": 10,
  "low": 20,
  "medium": 50,
  "high": 100
}
```