## 1. Ping Handling Flow
Whenever a ping result comes in:
1. Always store the result (success or failure).
2. Look at recent results to see how many **consecutive failures** there are.
3. Check whether an incident is already open for this monitor.
4. Based on these two things, decide whether to:
    - Start an incident,
    - Continue an existing incident,
    - Resolve an incident,
    - Or do nothing.
This is the core logic.
## 2. Incident Creation Flow (when a monitor fails)
When a ping fails:
1. Count how many failures happened in a row.
2. If the number of failures is **below the threshold**:
    - Do nothing special yet.
    - The system keeps waiting for more results.
3. If the failures reach the **threshold**:
    - Check if there is already an open incident.
    - If **no incident is open**:
        - Create a new incident.
        - Record an event saying the incident started.
        - Send notifications to the defined channels.
        - Record another event saying notifications were sent.
    - If **an incident is already open**:
        - Add an event saying more failures occurred.
        - Optionally send a reminder notification depending on your policy.
This ensures that exactly one active incident exists per monitor.
## 3. Incident Update Flow (when failures continue)
1. When a new failed ping arrives during an open incident, first check the **last recorded failure type** for that incident.
2. If the new failure has the **same type** as the previous one, do not add a new event.
3. If the new failure has a **different type**, add a new event to record this change.
4. This keeps the timeline clean by recording only meaningful changes while still reflecting how the situation evolves.
## 4. Incident Recovery Flow (when a ping succeeds)
When a ping succeeds:
1. Count how many **consecutive successes** occurred after the failure streak.
2. If no incident is open:
    - The monitor is simply healthy again.
    - No special action is needed.
3. If an incident _is_ open:
    - If the number of successes is **below the recovery threshold**:
        - Add an event such as “still checking, partial recovery”.
        - The incident stays open.
    - If the number of successes reaches the **recovery threshold**:
        - Mark the incident as resolved.
        - Add an event stating it recovered.
        - Send a “resolved” notification.
        - End the incident cleanly.