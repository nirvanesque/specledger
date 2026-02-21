# PostgREST API Contracts: sl revise

All requests require:
```
Authorization: Bearer <access_token>
apikey: <supabase_anon_key>
```

## 1. Get Project by Repo

```
GET /rest/v1/projects?repo_owner=eq.{owner}&repo_name=eq.{name}&select=id,default_branch
```

**Response**: `[{"id": "uuid", "default_branch": "main"}]`

## 2. Get Spec by Key

```
GET /rest/v1/specs?project_id=eq.{project_id}&spec_key=eq.{spec_key}&select=id,spec_key,phase
```

**Response**: `[{"id": "uuid", "spec_key": "136-revise-comments", "phase": "plan"}]`

## 3. Get Change for Spec

```
GET /rest/v1/changes?spec_id=eq.{spec_id}&select=id,head_branch,base_branch,state
```

**Response**: `[{"id": "uuid", "head_branch": "136-revise-comments", "base_branch": "main", "state": "open"}]`

## 4. Fetch Unresolved Comments

```
GET /rest/v1/review_comments?change_id=eq.{change_id}&is_resolved=eq.false&parent_comment_id=is.null&select=id,file_path,content,selected_text,line,start_line,author_name,author_email,created_at&order=created_at.asc
```

**Notes**:
- `parent_comment_id=is.null` filters to top-level comments only (excludes threaded replies)
- `order=created_at.asc` shows oldest comments first

**Response**:
```json
[{
  "id": "uuid",
  "file_path": "specledger/006-xxx/spec.md",
  "content": "this is unclear...",
  "selected_text": "when artifact content fails...",
  "line": null,
  "start_line": null,
  "author_name": "so0k",
  "author_email": "user@example.com",
  "created_at": "2026-02-19T12:42:15.286686+00:00"
}]
```

## 5. Resolve a Comment

```
PATCH /rest/v1/review_comments?id=eq.{comment_uuid}
Content-Type: application/json
Prefer: return=minimal

{"is_resolved": true}
```

**Response**: `204 No Content` on success

## 6. List Specs with Unresolved Comments (for branch picker)

No single PostgREST endpoint. Client-side aggregation:

```
GET /rest/v1/review_comments?is_resolved=eq.false&select=id,change_id,changes!inner(spec_id,specs!inner(spec_key,project_id))
```

Alternative (simpler, 2 calls):
1. Get all specs for project: `GET /rest/v1/specs?project_id=eq.{pid}&select=id,spec_key`
2. Get all changes for those specs: done client-side
3. Count unresolved comments per change: done client-side

Preferred approach: Fetch all unresolved comments for the project's changes and group client-side.

## Error Responses

| Status | Meaning | Action |
|--------|---------|--------|
| 401 | JWT expired/invalid | Prompt `sl auth login` |
| 403 | Not a project member | Check access |
| 404 | Resource not found | Verify spec_key/branch |
| PGRST303 | JWT expired | Prompt `sl auth login` |
