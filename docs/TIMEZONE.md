# Timezone Configuration

## Overview

The Simple Inventory API uses **UTC (Coordinated Universal Time)** for all timestamps throughout the application. This ensures consistency across different deployments, timezones, and client locations.

## Implementation

### 1. Application-Level UTC Enforcement

In [cmd/api/main.go:15](../cmd/api/main.go#L15):
```go
time.Local = time.UTC
```

This sets the default timezone for the entire Go application to UTC. All `time.Now()` calls will return UTC time.

### 2. Database-Level UTC Enforcement

In [internal/infrastructure/database/database.go:20-22](../internal/infrastructure/database/database.go#L20-L22):
```go
NowFunc: func() time.Time {
    return time.Now().UTC()
}
```

This ensures GORM uses UTC for all database timestamp operations, including `CreatedAt` and `UpdatedAt` fields.

### 3. Response Formatting

All API responses format timestamps in **ISO 8601 / RFC3339** format with UTC timezone.

Example: `2024-11-22T10:30:45Z`

The utility function in [internal/interfaces/http/util/time_util.go](../internal/interfaces/http/util/time_util.go):
```go
func FormatTimeUTC(t time.Time) string {
    return t.UTC().Format(time.RFC3339)
}
```

This is used consistently across all handlers:
- [AuthHandler](../internal/interfaces/http/handler/auth_handler.go#L60)
- [ProductHandler](../internal/interfaces/http/handler/product_handler.go#L211-L212)
- [LocationHandler](../internal/interfaces/http/handler/location_handler.go#L185-L186)
- [InventoryHandler](../internal/interfaces/http/handler/inventory_handler.go#L158)

## Benefits

### 1. **Consistency**
All timestamps are in the same timezone regardless of where the server is deployed.

### 2. **No Ambiguity**
UTC has no daylight saving time changes, eliminating edge cases around time transitions.

### 3. **Global Compatibility**
Clients in any timezone can easily convert UTC to their local time.

### 4. **Database Portability**
Moving databases between servers in different timezones doesn't affect data integrity.

### 5. **Simplified Development**
Developers don't need to worry about timezone conversions in business logic.

## Client-Side Handling

Clients should:

1. **Parse UTC timestamps** from API responses
2. **Convert to local timezone** for display
3. **Send timestamps in UTC** when creating/updating records (optional)

### JavaScript Example

```javascript
// Parse UTC timestamp from API
const utcTimestamp = "2024-11-22T10:30:45Z";
const date = new Date(utcTimestamp);

// Display in user's local timezone
console.log(date.toLocaleString()); // Converts to user's timezone

// Display in specific timezone
console.log(date.toLocaleString('en-US', { timeZone: 'America/New_York' }));
```

### Go Client Example

```go
// Parse UTC timestamp from API
timestamp, _ := time.Parse(time.RFC3339, "2024-11-22T10:30:45Z")

// Convert to specific timezone
loc, _ := time.LoadLocation("America/New_York")
localTime := timestamp.In(loc)
```

## Database Storage

PostgreSQL stores timestamps with timezone awareness:

```sql
-- Timestamps are stored as TIMESTAMPTZ (timestamp with timezone)
-- GORM automatically uses TIMESTAMPTZ for time.Time fields

-- Example query to view timestamps in different timezones
SELECT
    created_at AT TIME ZONE 'UTC' as utc_time,
    created_at AT TIME ZONE 'America/New_York' as ny_time,
    created_at AT TIME ZONE 'Asia/Tokyo' as tokyo_time
FROM products;
```

## API Response Examples

### Login Response
```json
{
  "token": "abc123...",
  "expires_at": "2024-11-23T10:30:45Z",
  "user": {
    "id": 1,
    "username": "admin"
  }
}
```

### Product Response
```json
{
  "id": 1,
  "name": "Laptop",
  "sku": "LAP-001",
  "created_at": "2024-11-22T10:30:45Z",
  "updated_at": "2024-11-22T15:20:30Z"
}
```

### Transaction Response
```json
{
  "id": 1,
  "product_id": 5,
  "type": "IN",
  "quantity": 50,
  "created_at": "2024-11-22T10:30:45Z"
}
```

## Testing Timezone Behavior

### Verify Application Timezone
```bash
# Run the API and check logs
go run cmd/api/main.go

# You should see:
# Timezone set to UTC
# Database timezone set to UTC
```

### Test API Responses
```bash
# Create a product and check the timestamp format
curl -X POST http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Test","sku":"TEST-001","price":10}'

# Response will include:
# "created_at": "2024-11-22T10:30:45Z"  <-- Note the 'Z' indicating UTC
```

### Database Verification
```sql
-- Connect to PostgreSQL and check timezone
SHOW timezone;  -- Should show 'UTC' or your server's timezone

-- Check how timestamps are stored
SELECT created_at,
       created_at AT TIME ZONE 'UTC' as explicit_utc
FROM products
LIMIT 1;
```

## Troubleshooting

### Issue: Timestamps appear in wrong timezone

**Cause**: Client not converting UTC to local time

**Solution**: Ensure client code converts UTC timestamps to local timezone for display

### Issue: Database shows different timezone

**Cause**: PostgreSQL server timezone setting

**Solution**: The application forces UTC usage via GORM's `NowFunc`, so this shouldn't affect data

### Issue: Session expiration seems incorrect

**Cause**: Comparing UTC time with local time

**Solution**: All time comparisons in the app use UTC, ensure client also uses UTC for comparisons

## Best Practices

1. ✅ **Always store timestamps in UTC**
2. ✅ **Convert to local timezone only for display**
3. ✅ **Use RFC3339 format for API communication**
4. ✅ **Test across different timezones**
5. ❌ **Never store timestamps in local timezone**
6. ❌ **Don't perform timezone conversions in business logic**

## Configuration

No additional configuration is needed. UTC is enforced at application startup.

If you need to change this behavior in the future (not recommended), modify:
- `cmd/api/main.go` - Remove `time.Local = time.UTC`
- `internal/infrastructure/database/database.go` - Modify `NowFunc`
- All handlers - Update timestamp formatting

## References

- [Go time package](https://pkg.go.dev/time)
- [RFC3339 Specification](https://datatracker.ietf.org/doc/html/rfc3339)
- [PostgreSQL Timezone Handling](https://www.postgresql.org/docs/current/datatype-datetime.html)
- [ISO 8601 Standard](https://en.wikipedia.org/wiki/ISO_8601)
