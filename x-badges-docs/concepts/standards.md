# Standards

Standards are informational tags that provide guidance on how to interpret and implement collection features. The collection interface is very feature-rich, and oftentimes you may need certain features to be implemented in a certain way, avoid certain features, etc. That is what standards are for.

## Purpose

All collections implement the same interface on the blockchain, but standards define:

-   How specific fields should be interpreted
-   Which features should be used or avoided
-   Expected metadata formats
-   Implementation guidelines for applications

## Timeline Implementation

```json
"standardsTimeline": [
  {
    "timelineTimes": [{"start": "1", "end": "18446744073709551615"}],
    "standards": ["transferable", "text-only-metadata", "non-fungible", "attendance-format"]
  }
]
```

## Important Notes

-   **No blockchain validation** - Standards are purely informational
-   **Multiple standards allowed** - As long as they are compatible
-   **Application responsibility** - Queriers must verify compliance

## Example Usage

```json
{
    "standardsTimeline": [
        {
            "timelineTimes": [{ "start": "1", "end": "18446744073709551615" }],
            "standards": ["soulbound", "event-attendance", "minimal-metadata"]
        }
    ]
}
```

Standards provide flexible guidance for collection behavior while maintaining blockchain simplicity.
