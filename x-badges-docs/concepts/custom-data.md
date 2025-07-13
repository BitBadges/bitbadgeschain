# Custom Data

Custom data fields are generic string fields that allow you to store any arbitrary value within BitBadges structures. These fields provide flexibility for storing application-specific information. They are not used for any specific purpose via the BitBadges site and are more for future customization and extensibility.

## Overview

Custom data fields appear throughout BitBadges as generic string storage:

-   **`customData`** - Simple string field in various structures
-   **`customDataTimeline`** - Timeline-based custom data that can change over time
-   **Custom fields in messages** - Additional data in transaction messages

## Usage

### Simple Custom Data

```json
{
    "customData": "Any string value you want to store"
}
```

### Timeline-Based Custom Data

```json
"customDataTimeline": [
  {
    "timelineTimes": [{"start": "1", "end": "18446744073709551615"}],
    "customData": "Application-specific data that changes over time"
  }
]
```

## Where You'll Find Custom Data

Custom data fields appear in:

-   **Collections** - `customDataTimeline` for collection-level data
-   **Address Lists** - `customData` for list-specific information
-   **Badge Metadata** - `customData` within badge metadata structures
-   **Messages** - Various transaction messages include custom data fields
