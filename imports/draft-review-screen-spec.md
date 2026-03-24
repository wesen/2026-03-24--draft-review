# Draft Review — Complete Screen Specification
## Retro Mac OS 1 Article Beta Reading Platform

> Every screen described as: ASCII mockup → Functional description → YAML DSL
> using existing React components: `MacWindow`, `TitleBar`, `Btn`, `HatchBar`,
> `ProgressBar`, `SectionNav`, `Paragraph`, `ReactionBadge`, `MenuBar`

---

## TABLE OF CONTENTS

```
 AUTHENTICATION
  01. Login
  02. Sign Up
  03. Forgot Password
  04. Email Verification

 AUTHOR SCREENS
  05. Dashboard (home) .............. ✅ BUILT
  06. Article Manager
  07. New Article / Upload
  08. Article Editor (section mgmt)
  09. Article Settings & Sharing
  10. Invite Readers (modal) ........ ✅ BUILT
  11. Single Article Analytics
  12. Feedback Deep Dive
  13. Feedback Export
  14. Version History
  15. Reader Management

 READER SCREENS
  16. Reader Welcome / Landing
  17. Reader View ................... ✅ BUILT
  18. Reader Summary (post-read)
  19. Reader History (my reviews)

 ACCOUNT & SETTINGS
  20. Profile & Account
  21. Notification Preferences
  22. Billing & Subscription

 ADMIN / TEAM
  23. Team Management
  24. Admin Overview

 MISC
  25. 404 / Error
  26. Onboarding Tour
  27. Keyboard Shortcuts (modal)
```

---

## GLOBAL SHELL

Every screen lives inside the global shell: a Mac OS 1 desktop with
the menu bar across the top and the checked-pattern background.

```
┌──────────────────────────────────────────────────────────────────┐
│ 🍎  File  Edit  View  Help                        ◉ 3 reactions │
├──────────────────────────────────────────────────────────────────┤
│░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░│
│░░░ ┌─────────────────────────────────────────────────────┐ ░░░░░│
│░░░ │▒▒▒▒▒▒▒▒▒▒▒▒▒▒ Window Title ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒│ ░░░░░│
│░░░ ├─────────────────────────────────────────────────────┤ ░░░░░│
│░░░ │                                                     │ ░░░░░│
│░░░ │              [ SCREEN CONTENT ]                     │ ░░░░░│
│░░░ │                                                     │ ░░░░░│
│░░░ └─────────────────────────────────────────────────────┘ ░░░░░│
│░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░│
└──────────────────────────────────────────────────────────────────┘
```

```yaml
shell:
  component: MenuBar
  menus:
    apple: [About Draft Review, separator, Version 1.0]
    file: [New Article… ⌘N, Import from Docs… ⌘I, separator, Export Feedback… ⌘E, separator, Log Out]
    edit: [Undo ⌘Z, Redo ⌘⇧Z, separator, Cut ⌘X, Copy ⌘C, Paste ⌘V]
    view: [Dashboard, Articles, Readers, separator, Show Reactions, Compact Mode]
    help: [Beta Reading Guide, How to Invite Readers, separator, Keyboard Shortcuts, Contact Support]
  right_status: "{reactions_count} reactions · {reader_name || author_name}"
  desktop:
    background: "checked_pattern(#e8e8e8, #c0c0c0, 4px)"
```

---

## 01 · LOGIN

```
┌──────────────────────────────────────────────────────────────────┐
│ 🍎  Draft Review                                                │
├──────────────────────────────────────────────────────────────────┤
│░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░│
│░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░│
│░░░░░  ┌─────────────────────────────────────────┐  ░░░░░░░░░░░░░│
│░░░░░  │☐│▒▒▒▒▒▒▒▒▒▒ Welcome Back ▒▒▒▒▒▒▒▒▒▒▒▒▒│  ░░░░░░░░░░░░│
│░░░░░  ├─────────────────────────────────────────┤  ░░░░░░░░░░░░░│
│░░░░░  │                                         │  ░░░░░░░░░░░░░│
│░░░░░  │              📝                         │  ░░░░░░░░░░░░░│
│░░░░░  │         Draft Review                    │  ░░░░░░░░░░░░░│
│░░░░░  │     Gather better feedback.             │  ░░░░░░░░░░░░░│
│░░░░░  │                                         │  ░░░░░░░░░░░░░│
│░░░░░  │  Email                                  │  ░░░░░░░░░░░░░│
│░░░░░  │  ┌───────────────────────────────────┐  │  ░░░░░░░░░░░░░│
│░░░░░  │  │ you@example.com                   │  │  ░░░░░░░░░░░░░│
│░░░░░  │  └───────────────────────────────────┘  │  ░░░░░░░░░░░░░│
│░░░░░  │                                         │  ░░░░░░░░░░░░░│
│░░░░░  │  Password                               │  ░░░░░░░░░░░░░│
│░░░░░  │  ┌───────────────────────────────────┐  │  ░░░░░░░░░░░░░│
│░░░░░  │  │ ••••••••••                        │  │  ░░░░░░░░░░░░░│
│░░░░░  │  └───────────────────────────────────┘  │  ░░░░░░░░░░░░░│
│░░░░░  │                                         │  ░░░░░░░░░░░░░│
│░░░░░  │  ╔═══════════════════════════════════╗  │  ░░░░░░░░░░░░░│
│░░░░░  │  ║           Sign In                 ║  │  ░░░░░░░░░░░░░│
│░░░░░  │  ╚═══════════════════════════════════╝  │  ░░░░░░░░░░░░░│
│░░░░░  │                                         │  ░░░░░░░░░░░░░│
│░░░░░  │     Forgot password?   Sign Up →        │  ░░░░░░░░░░░░░│
│░░░░░  │                                         │  ░░░░░░░░░░░░░│
│░░░░░  │  ┄┄┄┄┄┄┄┄┄ or ┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄  │  ░░░░░░░░░░░░░│
│░░░░░  │                                         │  ░░░░░░░░░░░░░│
│░░░░░  │  ┌───────────────────────────────────┐  │  ░░░░░░░░░░░░░│
│░░░░░  │  │  Continue with Google              │  │  ░░░░░░░░░░░░░│
│░░░░░  │  └───────────────────────────────────┘  │  ░░░░░░░░░░░░░│
│░░░░░  │                                         │  ░░░░░░░░░░░░░│
│░░░░░  └─────────────────────────────────────────┘  ░░░░░░░░░░░░░│
│░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░│
└──────────────────────────────────────────────────────────────────┘
```

**Functionality:** Standard email + password login. Magic link option for
readers who received an invite link. Google OAuth as secondary option.
Error states show inline under the field with a ⚠ prefix. After 3
failed attempts, a CAPTCHA dialog appears in a separate MacWindow.

```yaml
screen: login
window:
  component: MacWindow
  title: "Welcome Back"
  centered: true
  size: [400, 480]
  draggable: true

layout:
  - icon: { emoji: "📝", size: 48, border: true }
  - heading: { text: "Draft Review", font: CHICAGO, size: 18 }
  - subheading: { text: "Gather better feedback.", color: P.dark, size: 11 }
  - spacer: 16

  - field:
      id: email
      label: "Email"
      type: email
      placeholder: "you@example.com"
      style: mac_input  # 2px black border, GENEVA font, white bg
      validation: { required: true, pattern: email }
      error_display: inline_below  # ⚠ icon + message

  - field:
      id: password
      label: "Password"
      type: password
      placeholder: "••••••••"
      style: mac_input

  - button:
      component: Btn
      label: "Sign In"
      primary: true
      full_width: true
      action: auth.login(email, password)
      loading_state: "Signing in…"

  - row:
      justify: space-between
      children:
        - link: { text: "Forgot password?", action: navigate(/forgot) }
        - link: { text: "Sign Up →", action: navigate(/signup) }

  - divider: { text: "or", style: dotted }

  - button:
      label: "Continue with Google"
      full_width: true
      icon: google_icon
      action: auth.google_oauth()

special_cases:
  - invite_token_present:
      # If URL contains ?invite=xxx, show "You've been invited to review..."
      # Pre-fill email if known, skip to reader view after login
      banner: "You've been invited to review a draft"
      auto_redirect: /read/{article_id}

  - already_authenticated:
      redirect: /dashboard
```

---

## 02 · SIGN UP

```
┌──────────────────────────────────────────────────────────────────┐
│ 🍎  Draft Review                                                │
├──────────────────────────────────────────────────────────────────┤
│░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░│
│░░░░  ┌──────────────────────────────────────────┐  ░░░░░░░░░░░░░│
│░░░░  │☐│▒▒▒▒▒▒▒▒▒▒▒ Create Account ▒▒▒▒▒▒▒▒▒▒▒│  ░░░░░░░░░░░░│
│░░░░  ├──────────────────────────────────────────┤  ░░░░░░░░░░░░░│
│░░░░  │                                          │  ░░░░░░░░░░░░░│
│░░░░  │  I am a…                                 │  ░░░░░░░░░░░░░│
│░░░░  │  ┌──────────────┐  ┌──────────────────┐  │  ░░░░░░░░░░░░░│
│░░░░  │  │ ✎ AUTHOR     │  │ ◉ READER         │  │  ░░░░░░░░░░░░░│
│░░░░  │  │              │  │                  │  │  ░░░░░░░░░░░░░│
│░░░░  │  │ I'm writing  │  │ I was invited    │  │  ░░░░░░░░░░░░░│
│░░░░  │  │ something &  │  │ to give feedback │  │  ░░░░░░░░░░░░░│
│░░░░  │  │ need feedback │  │ on a draft       │  │  ░░░░░░░░░░░░░│
│░░░░  │  └──────────────┘  └──────────────────┘  │  ░░░░░░░░░░░░░│
│░░░░  │                                          │  ░░░░░░░░░░░░░│
│░░░░  │  Name                                    │  ░░░░░░░░░░░░░│
│░░░░  │  ┌────────────────────────────────────┐  │  ░░░░░░░░░░░░░│
│░░░░  │  │ Your name                          │  │  ░░░░░░░░░░░░░│
│░░░░  │  └────────────────────────────────────┘  │  ░░░░░░░░░░░░░│
│░░░░  │  Email                                   │  ░░░░░░░░░░░░░│
│░░░░  │  ┌────────────────────────────────────┐  │  ░░░░░░░░░░░░░│
│░░░░  │  │ you@example.com                    │  │  ░░░░░░░░░░░░░│
│░░░░  │  └────────────────────────────────────┘  │  ░░░░░░░░░░░░░│
│░░░░  │  Password                                │  ░░░░░░░░░░░░░│
│░░░░  │  ┌────────────────────────────────────┐  │  ░░░░░░░░░░░░░│
│░░░░  │  │ ••••••••                           │  │  ░░░░░░░░░░░░░│
│░░░░  │  └────────────────────────────────────┘  │  ░░░░░░░░░░░░░│
│░░░░  │                                          │  ░░░░░░░░░░░░░│
│░░░░  │  ╔════════════════════════════════════╗  │  ░░░░░░░░░░░░░│
│░░░░  │  ║         Create Account             ║  │  ░░░░░░░░░░░░░│
│░░░░  │  ╚════════════════════════════════════╝  │  ░░░░░░░░░░░░░│
│░░░░  │                                          │  ░░░░░░░░░░░░░│
│░░░░  │     Already have an account? Sign In     │  ░░░░░░░░░░░░░│
│░░░░  └──────────────────────────────────────────┘  ░░░░░░░░░░░░░│
│░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░│
└──────────────────────────────────────────────────────────────────┘
```

**Functionality:** Role selection (Author vs Reader) determines the
post-signup flow. Authors go to onboarding tour → first article creation.
Readers who arrived via invite link have "Reader" pre-selected and
their email pre-filled. Reader accounts are lightweight — no billing,
no article management. An account can upgrade from Reader to Author later.

```yaml
screen: signup
window:
  component: MacWindow
  title: "Create Account"
  centered: true
  size: [420, 540]

layout:
  - role_selector:
      label: "I am a…"
      options:
        - id: author
          icon: "✎"
          title: "AUTHOR"
          desc: "I'm writing something & need feedback"
        - id: reader
          icon: "◉"
          title: "READER"
          desc: "I was invited to give feedback on a draft"
      style: two_column_cards  # 2px black border, selected=inverted
      default: author  # unless invite_token in URL → reader

  - field: { id: name, label: "Name", placeholder: "Your name" }
  - field: { id: email, label: "Email", type: email }
  - field: { id: password, label: "Password", type: password, min_length: 8 }

  - button:
      component: Btn
      label: "Create Account"
      primary: true
      full_width: true
      action: auth.signup(role, name, email, password)

  - link: { text: "Already have an account? Sign In", action: navigate(/login) }

post_signup:
  author: navigate(/onboarding)
  reader_with_invite: navigate(/read/{article_id})
  reader_without_invite: navigate(/reader-home)
```

---

## 03 · FORGOT PASSWORD

```
┌──────────────────────────────────────────────────────────────────┐
│ 🍎  Draft Review                                                │
├──────────────────────────────────────────────────────────────────┤
│░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░│
│░░░░░  ┌─────────────────────────────────────────┐  ░░░░░░░░░░░░░│
│░░░░░  │☐│▒▒▒▒▒▒▒▒▒▒ Reset Password ▒▒▒▒▒▒▒▒▒▒▒│  ░░░░░░░░░░░░│
│░░░░░  ├─────────────────────────────────────────┤  ░░░░░░░░░░░░░│
│░░░░░  │                                         │  ░░░░░░░░░░░░░│
│░░░░░  │  Enter your email and we'll send you    │  ░░░░░░░░░░░░░│
│░░░░░  │  a link to reset your password.         │  ░░░░░░░░░░░░░│
│░░░░░  │                                         │  ░░░░░░░░░░░░░│
│░░░░░  │  Email                                  │  ░░░░░░░░░░░░░│
│░░░░░  │  ┌───────────────────────────────────┐  │  ░░░░░░░░░░░░░│
│░░░░░  │  │ you@example.com                   │  │  ░░░░░░░░░░░░░│
│░░░░░  │  └───────────────────────────────────┘  │  ░░░░░░░░░░░░░│
│░░░░░  │                                         │  ░░░░░░░░░░░░░│
│░░░░░  │  ╔═══════════════════════════════════╗  │  ░░░░░░░░░░░░░│
│░░░░░  │  ║         Send Reset Link           ║  │  ░░░░░░░░░░░░░│
│░░░░░  │  ╚═══════════════════════════════════╝  │  ░░░░░░░░░░░░░│
│░░░░░  │                                         │  ░░░░░░░░░░░░░│
│░░░░░  │  ─ ─ ─ SUCCESS STATE ─ ─ ─ ─ ─ ─ ─ ─  │  ░░░░░░░░░░░░░│
│░░░░░  │                                         │  ░░░░░░░░░░░░░│
│░░░░░  │              ✉                          │  ░░░░░░░░░░░░░│
│░░░░░  │        Check your email!                │  ░░░░░░░░░░░░░│
│░░░░░  │  We sent a reset link to a@b.com.       │  ░░░░░░░░░░░░░│
│░░░░░  │  Didn't get it? [Resend]                │  ░░░░░░░░░░░░░│
│░░░░░  │                                         │  ░░░░░░░░░░░░░│
│░░░░░  │          ← Back to Sign In              │  ░░░░░░░░░░░░░│
│░░░░░  └─────────────────────────────────────────┘  ░░░░░░░░░░░░░│
│░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░│
└──────────────────────────────────────────────────────────────────┘
```

```yaml
screen: forgot_password
window:
  component: MacWindow
  title: "Reset Password"
  centered: true
  size: [400, 320]

states:
  form:
    layout:
      - text: "Enter your email and we'll send you a link to reset your password."
      - field: { id: email, label: "Email", type: email }
      - button: { label: "Send Reset Link", primary: true, action: auth.reset(email) }
      - link: { text: "← Back to Sign In", action: navigate(/login) }

  success:
    layout:
      - icon: { emoji: "✉", size: 32 }
      - heading: "Check your email!"
      - text: "We sent a reset link to {email}."
      - link: { text: "Didn't get it? Resend", action: auth.reset(email), cooldown: 60s }
      - link: { text: "← Back to Sign In", action: navigate(/login) }
```

---

## 04 · EMAIL VERIFICATION

```
┌─────────────────────────────────────────┐
│☐│▒▒▒▒▒▒▒▒ Verify Your Email ▒▒▒▒▒▒▒▒▒▒│
├─────────────────────────────────────────┤
│                                         │
│                 ✉                       │
│                                         │
│       We sent a verification link       │
│       to alex@example.com               │
│                                         │
│  ┌───────────────────────────────────┐  │
│  │  Enter 6-digit code:             │  │
│  │  ┌──┐ ┌──┐ ┌──┐ ┌──┐ ┌──┐ ┌──┐  │  │
│  │  │ 4│ │ 7│ │  │ │  │ │  │ │  │  │  │
│  │  └──┘ └──┘ └──┘ └──┘ └──┘ └──┘  │  │
│  └───────────────────────────────────┘  │
│                                         │
│  ╔═══════════════════════════════════╗  │
│  ║            Verify                 ║  │
│  ╚═══════════════════════════════════╝  │
│                                         │
│    Didn't get it?  [Resend]  (47s)      │
│                                         │
│  ┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄  │
│  Wrong email? [Change email]            │
└─────────────────────────────────────────┘
```

```yaml
screen: email_verification
window:
  component: MacWindow
  title: "Verify Your Email"
  centered: true
  size: [400, 360]

layout:
  - icon: { emoji: "✉", size: 32 }
  - text: "We sent a verification link to {email}"
  - code_input:
      digits: 6
      style: retro_boxes  # individual 2px border boxes, CHICAGO font
      auto_submit: true  # submits when all 6 entered
  - button: { label: "Verify", primary: true }
  - resend_link: { cooldown: 60s, text: "Didn't get it? Resend" }
  - divider: dotted
  - link: { text: "Wrong email? Change email", action: navigate(/signup) }
```

---

## 05 · DASHBOARD (AUTHOR HOME) — ✅ ALREADY BUILT

```
┌──────────────────────────────────────────────────────────────────┐
│ 🍎  File  Edit  View  Help                    alex@example.com  │
├──────────────────────────────────────────────────────────────────┤
│☐│▒▒▒▒▒▒▒▒▒▒ Draft Review — Dashboard ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒│
├──────────────────────────────────────────────────────────────────┤
│ [In Review] Why Design Systems Fail  [Draft] Remote Work  [New] │
│                                                                  │
│ ┌─────────────┐┌─────────────┐┌─────────────┐┌─────────────┐   │
│ │ ◉ 5         ││ ✦ 10        ││ § 5          ││ ▸ 60%       │   │
│ │ READERS     ││ REACTIONS   ││ SECTIONS    ││ AVG PROGRESS │   │
│ └─────────────┘└─────────────┘└─────────────┘└─────────────┘   │
│                                                                  │
│ ┌─── Readers ─────────┐  ┌─── Reactions by Section ──────────┐  │
│ │ SK Sarah K. ▓▓▓ 100%│  │ Introduction    ██░░  3 reactions │  │
│ │ MT Marcus   ▓▓░  60%│  │ Adoption Cliff  ████  4 reactions │  │
│ │ PR Priya R. ▓▓▓  80%│  │ Governance      ███░  2 reactions │  │
│ │ JL James L. ▓░░  40%│  │ Documentation   ██░░  2 reactions │  │
│ │ CW Chen W.  ▓░░  20%│  │ Making It Work  █░░░  1 reaction  │  │
│ │         [+ Invite]  │  │ ★ Useful  ? Confused  ◎ Slow  ♥  │  │
│ └─────────────────────┘  └───────────────────────────────────┘  │
│                                                                  │
│ ┌ ⚠ Potential Draft-Killer Detected ─────────────────────────┐  │
│ │ "The Adoption Cliff" has the most confusion and pacing     │  │
│ │ issues (3 flags). Focus your next revision here.           │  │
│ │                            [Review This Section →]         │  │
│ └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│ ┌─── Recent Feedback ────────────────────────────────────────┐  │
│ │ ★ Sarah K. on Introduction — Great hook                    │  │
│ │ ? Marcus T. on Adoption Cliff — What do you mean...        │  │
│ │ ♥ Sarah K. on Adoption Cliff — YES. This is exactly...     │  │
│ └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│                  ╔══════════════════════╗                        │
│                  ║  Open Full Review →  ║                        │
│                  ╚══════════════════════╝                        │
└──────────────────────────────────────────────────────────────────┘
```

> **Implemented in:** `review-system.jsx` — See existing code.

---

## 06 · ARTICLE MANAGER

```
┌──────────────────────────────────────────────────────────────────┐
│☐│▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒ My Articles ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒│
├──────────────────────────────────────────────────────────────────┤
│  [+ New Article]    Search: [________________]    Sort: [Recent] │
│─────────────────────────────────────────────────────────────────│
│                                                                  │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │ ✎ Why Design Systems Fail                     [In Review] │  │
│  │   5 sections · 10 reactions · 5 readers                   │  │
│  │   v3 — Updated 2 days ago                                 │  │
│  │   ▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░ 60% avg read                        │  │
│  │   ⚠ Draft-killer in "Adoption Cliff"                      │  │
│  │                 [Analytics]  [Edit]  [Share]  [Archive]    │  │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │ ✎ Remote Work Isn't Working                       [Draft] │  │
│  │   3 sections · 2 reactions · 2 readers                    │  │
│  │   v1 — Updated 1 week ago                                 │  │
│  │   ▓▓▓▓▓▓▓▓░░░░░░░░░ 66% avg read                        │  │
│  │                 [Analytics]  [Edit]  [Share]  [Archive]    │  │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │ ✎ The Art of Saying No                             [New]  │  │
│  │   1 section · 0 reactions · 0 readers                     │  │
│  │   v1 — Created 3 days ago                                 │  │
│  │   No readers yet — invite your first beta reader!         │  │
│  │                 [Analytics]  [Edit]  [Share]  [Invite →]  │  │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┄┄┄┄┄┄┄ Archived (1) ┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄ [show] ┄┄┄┄┄┄┄┄  │
└──────────────────────────────────────────────────────────────────┘
```

**Functionality:** Master list of all articles. Each card shows key stats
at a glance: version number, reader count, reaction count, average
read-through, and draft-killer alerts. Actions: open analytics, edit
sections, manage sharing/invitations, archive. Search filters by title.
Sort by recent, most reactions, most readers, or status. Archived
articles collapse to a single expandable row.

```yaml
screen: article_manager
window:
  component: MacWindow
  title: "My Articles"
  maximized: true

toolbar:
  - button: { label: "+ New Article", action: navigate(/new-article) }
  - spacer: flex
  - search: { placeholder: "Search…", field: title, style: mac_input }
  - dropdown:
      label: "Sort"
      options: [Recent, Most Reactions, Most Readers, Status]
      style: mac_dropdown  # 2px border, CHICAGO font

list:
  component: scrollable_list
  item_template:
    component: article_card
    border: "2px solid P.black"
    layout:
      - row:
          - icon: "✎"
          - title: "{article.title}"
          - badge: { text: "{article.status}", style: status_pill }
      - text: "{sections} sections · {reactions} reactions · {readers} readers"
      - text: "v{version} — Updated {relative_time}"
      - component: HatchBar
        percent: "{avg_progress}"
        label: "{avg_progress}% avg read"
      - conditional:
          if: "book_killer != null"
          show: "⚠ Draft-killer in \"{book_killer.title}\""
      - action_row:
          buttons:
            - { label: "Analytics", action: navigate(/article/{id}/analytics) }
            - { label: "Edit", action: navigate(/article/{id}/edit) }
            - { label: "Share", action: navigate(/article/{id}/share) }
            - { label: "Archive", action: article.archive(id), confirm: true }

  empty_state:
    icon: "📝"
    heading: "No articles yet"
    text: "Create your first article and start getting feedback."
    button: { label: "New Article →", primary: true }

footer:
  - collapsible:
      label: "Archived ({count})"
      initially: collapsed
      content: same article_card but greyed, with [Restore] action
```

---

## 07 · NEW ARTICLE / UPLOAD

```
┌──────────────────────────────────────────────────────────────────┐
│☐│▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒ New Article ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒│
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Title                                                           │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │ Untitled Article                                          │  │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  Author Note (shown to readers before they begin)                │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │ Thanks for reading an early draft! I'm especially         │  │
│  │ looking for feedback on what's useful and where things     │  │
│  │ get confusing. Don't worry about typos.                    │  │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┄┄┄┄┄┄┄┄┄┄┄┄ Import Content ┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄  │
│                                                                  │
│  ┌──────────┐  ┌──────────┐  ┌───────────┐  ┌──────────────┐   │
│  │  📄      │  │  📝      │  │  📋       │  │    ⌨        │   │
│  │  .docx   │  │  Google  │  │  .md      │  │   Paste     │   │
│  │  upload  │  │  Docs    │  │  Markdown │  │   text      │   │
│  └──────────┘  └──────────┘  └───────────┘  └──────────────┘   │
│                                                                  │
│  ┌─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ┐  │
│  │                                                           │  │
│  │         Drop a .docx or .md file here                     │  │
│  │                or click to browse                         │  │
│  │                                                           │  │
│  └─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ┘  │
│                                                                  │
│  ┄ After upload, we'll split into sections by heading ┄         │
│                                                                  │
│                    [Cancel]  ╔══════════════╗                    │
│                              ║  Continue →  ║                    │
│                              ╚══════════════╝                    │
└──────────────────────────────────────────────────────────────────┘
```

**Functionality:** Create an article by uploading .docx, linking a Google
Doc, uploading Markdown, or pasting raw text. The system auto-splits
content into sections by heading (H1/H2). After import, user is taken
to the Article Editor to review and rearrange sections. The "Author
Note" is what readers see on the Welcome screen before they begin.
Google Docs import uses a read-only link and syncs content.

```yaml
screen: new_article
window:
  component: MacWindow
  title: "New Article"
  centered: true
  size: [560, 580]

layout:
  - field:
      id: title
      label: "Title"
      placeholder: "Untitled Article"
      style: mac_input
      font: CHICAGO
      size: 14

  - field:
      id: author_note
      label: "Author Note (shown to readers before they begin)"
      type: textarea
      rows: 3
      default: "Thanks for reading an early draft! I'm especially looking for feedback on what's useful and where things get confusing. Don't worry about typos."

  - divider: { text: "Import Content" }

  - import_methods:
      style: four_column_cards  # retro icon cards, 2px border
      options:
        - id: docx
          icon: "📄"
          label: ".docx upload"
          action: file_picker(accept=".docx")
        - id: gdocs
          icon: "📝"
          label: "Google Docs"
          action: google_docs_link_modal
        - id: markdown
          icon: "📋"
          label: "Markdown"
          action: file_picker(accept=".md")
        - id: paste
          icon: "⌨"
          label: "Paste text"
          action: show_paste_textarea

  - drop_zone:
      accept: [.docx, .md, .txt]
      text: "Drop a .docx or .md file here\nor click to browse"
      style: dashed_border
      on_drop: parse_and_preview

  - hint: "After upload, we'll split into sections by heading"

  - action_row:
      - button: { label: "Cancel", action: navigate(-1) }
      - button: { label: "Continue →", primary: true, disabled_until: content_loaded }

post_import:
  action: navigate(/article/{new_id}/edit)
  auto_split: by_heading  # H1, H2 become section breaks
```

---

## 08 · ARTICLE EDITOR (SECTION MANAGEMENT)

```
┌──────────────────────────────────────────────────────────────────┐
│☐│▒▒▒▒▒▒▒ Edit: Why Design Systems Fail ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒│
├──────────────────────────────────────────────────────────────────┤
│ ┌─ Sections ──────┐  ┌─ Section Content ──────────────────────┐ │
│ │                  │  │                                       │ │
│ │  ≡ 1. Intro     │  │  Title                                │ │
│ │  ≡ 2. Adoption  │  │  ┌─────────────────────────────────┐  │ │
│ │  ≡ 3. Governanc │  │  │ The Adoption Cliff              │  │ │
│ │  ≡ 4. Documenta │  │  └─────────────────────────────────┘  │ │
│ │  ≡ 5. Making It │  │                                       │ │
│ │                  │  │  Content                              │ │
│ │  drag to reorder │  │  ┌─────────────────────────────────┐  │ │
│ │                  │  │  │ Most design systems see strong  │  │ │
│ │ ┌──────────────┐ │  │  │ initial adoption. Engineers are │  │ │
│ │ │ + Add Section │ │  │  │ excited, designers feel heard, │  │ │
│ │ └──────────────┘ │  │  │ and leadership is optimistic.  │  │ │
│ │                  │  │  │                                │  │ │
│ │ ── Danger Zone ─ │  │  │ But around month six, usage    │  │ │
│ │ [Delete Section] │  │  │ plateaus. Teams start creating │  │ │
│ │                  │  │  │ workarounds, one-off components│  │ │
│ └──────────────────┘  │  │ multiply...                    │  │ │
│                       │  └─────────────────────────────────┘  │ │
│                       │                                       │ │
│                       │  Paragraph breaks become separate     │ │
│                       │  "paragraphs" in the reader view,     │ │
│                       │  each individually reactable.         │ │
│                       │                                       │ │
│                       │  [Preview as Reader]  ╔════════════╗  │ │
│                       │                       ║   Save     ║  │ │
│                       │                       ╚════════════╝  │ │
│                       └───────────────────────────────────────┘ │
└──────────────────────────────────────────────────────────────────┘
```

**Functionality:** Two-panel editor. Left sidebar lists sections with
drag handles (≡) for reordering. Click a section to edit its title and
content in the right pane. Content is plain text — paragraph breaks
(double newline) define the individual paragraphs that readers will
react to. "Preview as Reader" opens the Reader View in a new window.
"Add Section" appends a blank section. "Delete Section" has a
confirmation dialog. Saving creates a new version silently. A
"re-import" button in the toolbar allows re-uploading a new draft
(which creates a new version and preserves all old feedback).

```yaml
screen: article_editor
window:
  component: MacWindow
  title: "Edit: {article.title}"
  maximized: true

layout: split_pane
  left:
    width: 200
    border_right: "2px solid P.black"
    children:
      - header: { text: "Sections", font: CHICAGO, size: 10 }
      - sortable_list:
          items: article.sections
          drag_handle: "≡"
          on_reorder: article.reorder_sections
          item_template:
            text: "{index}. {section.title}"
            selected_style: inverted  # black bg, white text
            on_click: set_active_section(section.id)
      - spacer: flex
      - button: { label: "+ Add Section", full_width: true, action: article.add_section }
      - divider: { text: "Danger Zone", color: P.dark }
      - button:
          label: "Delete Section"
          full_width: true
          style: danger  # red-tinted or warning pattern
          action: confirm_then(article.delete_section)

  right:
    padding: 20
    children:
      - field:
          id: section_title
          label: "Title"
          value: "{active_section.title}"
          font: CHICAGO
          size: 14
      - field:
          id: section_content
          label: "Content"
          type: textarea
          value: "{active_section.paragraphs.join('\\n\\n')}"
          rows: auto_grow
          hint: "Paragraph breaks become separate reactable paragraphs in reader view."
          monospace: false
          font: GENEVA
      - action_row:
          - button: { label: "Preview as Reader", action: open_reader_preview }
          - button: { label: "Save", primary: true, action: article.save }

toolbar_extras:
  - button: { label: "Re-import…", action: show_reimport_modal }
  - version_badge: "v{article.version}"
```

---

## 09 · ARTICLE SETTINGS & SHARING

```
┌──────────────────────────────────────────────────────────────────┐
│☐│▒▒▒▒▒▒▒▒▒▒▒ Article Settings ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒│
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  SHARING ─────────────────────────────────────────────────────   │
│                                                                  │
│  Review Link                                                     │
│  ┌──────────────────────────────────────────┐                    │
│  │ https://draftreview.app/r/abc123         │  [Copy]  [Reset]  │
│  └──────────────────────────────────────────┘                    │
│                                                                  │
│  Access: ◉ Invite-only (readers need account)                    │
│          ○ Anyone with link (no login needed)                    │
│          ○ Password protected: [________]                        │
│                                                                  │
│  Reader visibility:                                              │
│  [✓] Readers can see each other's reactions                      │
│  [ ] Readers can see each other's names                          │
│  [✓] Show author note on welcome screen                          │
│                                                                  │
│  FEEDBACK ────────────────────────────────────────────────────   │
│                                                                  │
│  Enabled reactions:                                              │
│  [✓] ★ Useful    [✓] ? Confused                                 │
│  [✓] ◎ Slow      [✓] ♥ Favorite                                 │
│  [ ] 💡 Suggestion (custom)                                      │
│  [+ Add custom reaction]                                         │
│                                                                  │
│  [ ] Require a note with each reaction                           │
│  [✓] Allow anonymous reactions                                   │
│                                                                  │
│  STATUS ──────────────────────────────────────────────────────   │
│                                                                  │
│  ○ Draft   ◉ In Review   ○ Complete   ○ Archived                 │
│                                                                  │
│                              ╔════════════════╗                  │
│                              ║  Save Settings ║                  │
│                              ╚════════════════╝                  │
└──────────────────────────────────────────────────────────────────┘
```

**Functionality:** Controls how the article is shared and what feedback
options are available. The review link can be copied or reset (which
invalidates old links). Three access modes: invite-only, open link,
or password-protected. Reader visibility toggles control whether
readers can see other people's reactions (useful for preventing
groupthink in early rounds, per the beta reading guide). Authors can
enable/disable individual reaction types or add custom ones. Status
controls whether the article appears as active or archived.

```yaml
screen: article_settings
window:
  component: MacWindow
  title: "Article Settings"
  size: [520, 600]

layout:
  - section_header: "SHARING"

  - field_with_actions:
      id: review_link
      label: "Review Link"
      value: "{base_url}/r/{article.share_token}"
      readonly: true
      actions:
        - { label: "Copy", action: clipboard.copy }
        - { label: "Reset", action: article.reset_share_token, confirm: "This will break existing links." }

  - radio_group:
      id: access_mode
      label: "Access"
      options:
        - { value: invite_only, label: "Invite-only (readers need account)" }
        - { value: link, label: "Anyone with link (no login needed)" }
        - { value: password, label: "Password protected", sub_field: { id: password, type: password } }

  - checkbox_group:
      label: "Reader visibility"
      items:
        - { id: see_reactions, label: "Readers can see each other's reactions", default: true }
        - { id: see_names, label: "Readers can see each other's names", default: false }
        - { id: show_author_note, label: "Show author note on welcome screen", default: true }

  - section_header: "FEEDBACK"

  - reaction_toggles:
      label: "Enabled reactions"
      items:
        - { type: useful, icon: "★", default: true }
        - { type: confused, icon: "?", default: true }
        - { type: slow, icon: "◎", default: true }
        - { type: favorite, icon: "♥", default: true }
      custom_reactions:
        allow_add: true
        max: 3

  - checkbox_group:
      items:
        - { id: require_note, label: "Require a note with each reaction", default: false }
        - { id: allow_anon, label: "Allow anonymous reactions", default: true }

  - section_header: "STATUS"

  - radio_group:
      id: status
      options: [Draft, In Review, Complete, Archived]

  - button: { label: "Save Settings", primary: true, full_width: false }
```

---

## 10 · INVITE READERS — ✅ ALREADY BUILT (expanded spec)

> Already built as `InviteDialog`. Expanded here with batch invite.

```
┌─────────────────────────────────────────┐
│☐│▒▒▒▒▒▒▒▒▒ Invite Readers ▒▒▒▒▒▒▒▒▒▒▒▒│
├─────────────────────────────────────────┤
│                                         │
│  ┌──────────────┐ ┌─────────────────┐   │
│  │ ◉ Personal   │ │ ○ Batch Invite  │   │
│  └──────────────┘ └─────────────────┘   │
│                                         │
│  Email Address                          │
│  ┌───────────────────────────────────┐  │
│  │ reader@example.com                │  │
│  └───────────────────────────────────┘  │
│                                         │
│  Personal Note                          │
│  ┌───────────────────────────────────┐  │
│  │ I'm working on an article and     │  │
│  │ would love your honest feedback   │  │
│  │ — especially where you feel       │  │
│  │ confused or where things drag.    │  │
│  └───────────────────────────────────┘  │
│                                         │
│  ┄┄┄ Previously Invited ┄┄┄┄┄┄┄┄┄┄┄┄  │
│  ✓ sarah@ex.com — opened, 5 reactions  │
│  ✓ marcus@ex.com — opened, 2 reactions │
│  ◯ james@ex.com — not yet opened       │
│  ◯ chen@ex.com — not yet opened        │
│                                         │
│          [Cancel]  ╔══════════════╗     │
│                    ║ Send Invite  ║     │
│                    ╚══════════════╝     │
└─────────────────────────────────────────┘
```

```yaml
screen: invite_readers
type: modal
component: MacWindow  # dialog variant
title: "Invite Readers"
size: [440, 520]

tabs:
  personal:
    label: "Personal"
    layout:
      - field: { id: email, label: "Email Address", type: email }
      - field:
          id: note
          label: "Personal Note"
          type: textarea
          rows: 4
          default: "I'm working on an article and would love your honest feedback — especially where you feel confused or where things drag."

  batch:
    label: "Batch Invite"
    layout:
      - field:
          id: emails
          label: "Email Addresses (one per line)"
          type: textarea
          rows: 6
          placeholder: "reader1@example.com\nreader2@example.com"
      - field: { id: note, label: "Shared Note", type: textarea, rows: 3 }

previously_invited:
  label: "Previously Invited"
  style: collapsible_list
  item_template:
    - icon: "✓ if opened else ◯"
    - email
    - status: "{reactions_count} reactions" or "not yet opened"
    - action: [Resend] if not opened

actions:
  - button: { label: "Cancel" }
  - button: { label: "Send Invite", primary: true }

post_send:
  show_reminder: "Remember, only ~1 in 4 invitees will engage — that's normal!"
```

---

## 11 · SINGLE ARTICLE ANALYTICS

```
┌──────────────────────────────────────────────────────────────────┐
│☐│▒▒▒▒▒▒ Analytics: Why Design Systems Fail ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒│
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  READER PROGRESS ─────────────────────────────────────────────   │
│                                                                  │
│  s1 Intro       │████████████████████████████████████│ 5/5       │
│  s2 Adoption    │██████████████████████████░░░░░░░░░░│ 3/5       │
│  s3 Governance  │████████████████████░░░░░░░░░░░░░░░░│ 3/5       │
│  s4 Docs        │████████████████░░░░░░░░░░░░░░░░░░░░│ 2/5       │
│  s5 Making It   │████████████░░░░░░░░░░░░░░░░░░░░░░░░│ 1/5       │
│                                                 ↑ drop-off       │
│                                                                  │
│  REACTION HEATMAP ────────────────────────────────────────────   │
│                                                                  │
│        ★useful  ?confused  ◎slow  ♥favorite                      │
│  s1    ██       ░░         ░░     █                              │
│  s2    █        ██         █      ██                             │
│  s3    ██       █          ░░     ░░                             │
│  s4    █        ░░         ░░     █                              │
│  s5    ░░       ░░         █      ░░                             │
│                                                                  │
│  KEY INSIGHTS ────────────────────────────────────────────────   │
│                                                                  │
│  ⚠ DRAFT-KILLER: Readers drop off after "The Adoption Cliff"    │
│    3 of 5 readers stopped reading at or before Section 3.        │
│    This section has the highest confusion rate (2 flags).        │
│                                                                  │
│  ✓ WORKING WELL: "Introduction" — all readers passed through,   │
│    3 positive reactions, strong hook.                             │
│                                                                  │
│  ◎ NEEDS ATTENTION: "Making It Work" — only 1 reader reached    │
│    it, flagged as slow. May need more depth once readers         │
│    actually get here.                                            │
│                                                                  │
│  ┄┄ VERSION COMPARISON ┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄   │
│                                                                  │
│  v1 (Feb 12) → v2 (Feb 28) → v3 (Mar 10)                       │
│  Avg read:  32%      48%         60%      ↑ improving            │
│  Reactions:  4        7          10                              │
│  Readers:    3        4           5                              │
│                                                                  │
│        [Export Report]  [View All Feedback →]                    │
└──────────────────────────────────────────────────────────────────┘
```

**Functionality:** Deep analytics for a single article. The reader progress
chart shows how many readers reached each section — the "drop-off
curve" is the single most important chart, directly from the Help This
Book model. Reaction heatmap shows density of each type per section.
Key Insights section auto-generates natural language summaries
identifying the draft-killer, what's working, and what needs attention.
Version comparison shows improvement across revisions (avg read-through,
total reactions, reader count). Export generates a PDF or Markdown report.

```yaml
screen: article_analytics
window:
  component: MacWindow
  title: "Analytics: {article.title}"
  maximized: true

layout:
  - section: "READER PROGRESS"
    component: dropoff_chart
    type: horizontal_bars
    style: hatched  # retro fill pattern
    data: sections.map(s => { label: s.title, reached: readers_who_reached(s), total: readers.length })
    annotation:
      type: arrow
      label: "↑ drop-off"
      position: steepest_decline

  - section: "REACTION HEATMAP"
    component: heatmap_grid
    rows: sections
    columns: reaction_types
    cell_style: block_density  # more blocks = more reactions
    legend: [★useful, ?confused, ◎slow, ♥favorite]

  - section: "KEY INSIGHTS"
    component: auto_insights
    rules:
      draft_killer:
        icon: "⚠"
        condition: "section with highest (confused + slow) AND steepest dropoff"
        template: "Readers drop off after \"{section}\". {stats}."
      working_well:
        icon: "✓"
        condition: "section with highest pass-through AND positive reactions"
        template: "\"{section}\" — {stats}."
      needs_attention:
        icon: "◎"
        condition: "section with low reach OR high slow count"
        template: "\"{section}\" — {stats}."

  - section: "VERSION COMPARISON"
    component: version_timeline
    metrics: [avg_read_percent, total_reactions, total_readers]
    data: article.versions
    trend_indicator: true  # ↑ improving or ↓ declining

  - action_row:
      - button: { label: "Export Report", action: export.analytics_report }
      - button: { label: "View All Feedback →", action: navigate(/article/{id}/feedback) }
```

---

## 12 · FEEDBACK DEEP DIVE

```
┌──────────────────────────────────────────────────────────────────┐
│☐│▒▒▒▒▒▒▒▒▒ Feedback: Why Design Systems Fail ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒│
├──────────────────────────────────────────────────────────────────┤
│ Filters:                                                         │
│ Section: [All ▾]  Type: [All ▾]  Reader: [All ▾]  [Clear all]  │
│──────────────────────────────────────────────────────────────────│
│                                                                  │
│ ── s1: Introduction ──────────────────────────────────────────   │
│                                                                  │
│ ┌ ★ ─────────────────────────────────────────────────────────┐  │
│ │ Sarah K.  ·  2 days ago                        [Resolved ▾]│  │
│ │ "Great hook — immediately relevant to my situation."       │  │
│ │                                                            │  │
│ │ Author note: _____________________________________________ │  │
│ └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│ ┌ ★ ─────────────────────────────────────────────────────────┐  │
│ │ Chen W.  ·  3 days ago                          [New     ▾]│  │
│ │ "Strong opening. Sets the right expectations."             │  │
│ └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│ ── s2: The Adoption Cliff ────────────────────────────────────   │
│                                                                  │
│ ┌ ? ─────────────────────────────────────────────────────────┐  │
│ │ Marcus T.  ·  2 days ago                        [New     ▾]│  │
│ │ "What do you mean by 'workarounds'? Can you give..."       │  │
│ │                                                            │  │
│ │ Author note: [Added examples in v4 revision___] ✓ Saved   │  │
│ └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│ ┌ ◎ ─────────────────────────────────────────────────────────┐  │
│ │ Priya R.  ·  3 days ago                     [Noted      ▾]│  │
│ │ "This section drags a bit. Maybe tighten the middle..."    │  │
│ └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│           Showing 10 of 10 reactions · [Export as CSV]           │
└──────────────────────────────────────────────────────────────────┘
```

**Functionality:** All feedback for an article in one scrollable view,
grouped by section. Each reaction card shows the type icon, reader
name, timestamp, their note, and a status dropdown (New / Noted /
Resolved). Authors can add private "Author notes" to each reaction
to track what they plan to do about it (visible only to them). Filter
by section, reaction type, or reader. The status system helps authors
triage feedback across revision cycles. Export as CSV for external
use. Importantly, this follows the principle that "readers spot
problems, you design solutions" — the author note field is where the
author records their planned fix.

```yaml
screen: feedback_deep_dive
window:
  component: MacWindow
  title: "Feedback: {article.title}"
  maximized: true

toolbar:
  filters:
    - dropdown: { id: section, label: "Section", options: [All, ...sections] }
    - dropdown: { id: type, label: "Type", options: [All, Useful, Confused, Slow, Favorite] }
    - dropdown: { id: reader, label: "Reader", options: [All, ...readers] }
    - link: { text: "Clear all", action: clear_filters }

list:
  group_by: section
  group_header:
    style: "divider with section title"
    format: "── s{index}: {section.title} ──"

  item_template:
    component: feedback_card
    border: "2px solid P.black"
    layout:
      - row:
          - reaction_icon: "{reaction.type}"  # large icon left edge
          - column:
              - row:
                  - text: "{reader.name}  ·  {relative_time}"
                  - spacer: flex
                  - dropdown:
                      id: status
                      options: [New, Noted, Resolved, Dismissed]
                      style: mac_dropdown_small
              - text: "\"{reaction.text}\""
              - conditional:
                  if: "status == Noted || status == Resolved"
                  show:
                    field:
                      id: author_note
                      placeholder: "What are you going to do about this?"
                      style: inline_editable
                      private: true  # only visible to author

footer:
  - text: "Showing {filtered_count} of {total_count} reactions"
  - button: { label: "Export as CSV", action: export.feedback_csv }
```

---

## 13 · FEEDBACK EXPORT (MODAL)

```
┌─────────────────────────────────────────┐
│☐│▒▒▒▒▒▒▒▒ Export Feedback ▒▒▒▒▒▒▒▒▒▒▒▒│
├─────────────────────────────────────────┤
│                                         │
│  Article: Why Design Systems Fail       │
│  Version: v3 · 10 reactions · 5 readers │
│                                         │
│  Format:                                │
│  ◉ CSV (spreadsheet-friendly)           │
│  ○ Markdown (human-readable)            │
│  ○ JSON (developer-friendly)            │
│                                         │
│  Include:                               │
│  [✓] Reaction type & text               │
│  [✓] Reader name                        │
│  [✓] Section & paragraph reference      │
│  [✓] Timestamp                          │
│  [ ] Author notes (private)             │
│  [ ] Reader progress data               │
│                                         │
│  Filter:                                │
│  [✓] All reactions                      │
│  [ ] Only unresolved                    │
│  [ ] Only confused + slow               │
│                                         │
│         [Cancel]  ╔══════════════╗      │
│                   ║   Export     ║      │
│                   ╚══════════════╝      │
└─────────────────────────────────────────┘
```

```yaml
screen: feedback_export
type: modal
component: MacWindow
title: "Export Feedback"
size: [400, 440]

layout:
  - info: "Article: {article.title}\nVersion: v{version} · {reactions} reactions · {readers} readers"

  - radio_group:
      id: format
      label: "Format"
      options:
        - { value: csv, label: "CSV (spreadsheet-friendly)" }
        - { value: md, label: "Markdown (human-readable)" }
        - { value: json, label: "JSON (developer-friendly)" }

  - checkbox_group:
      label: "Include"
      items:
        - { id: type_text, label: "Reaction type & text", default: true }
        - { id: reader_name, label: "Reader name", default: true }
        - { id: section_ref, label: "Section & paragraph reference", default: true }
        - { id: timestamp, label: "Timestamp", default: true }
        - { id: author_notes, label: "Author notes (private)", default: false }
        - { id: progress, label: "Reader progress data", default: false }

  - radio_group:
      id: filter
      label: "Filter"
      options:
        - { value: all, label: "All reactions" }
        - { value: unresolved, label: "Only unresolved" }
        - { value: pain, label: "Only confused + slow" }

  - action_row:
      - button: { label: "Cancel" }
      - button: { label: "Export", primary: true, action: export.download(format, options) }
```

---

## 14 · VERSION HISTORY

```
┌──────────────────────────────────────────────────────────────────┐
│☐│▒▒▒▒▒▒▒▒▒ Version History ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒│
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Why Design Systems Fail                                         │
│                                                                  │
│  ● v3 (current) ─── Mar 10, 2026 ──────────────────────────    │
│  │  5 sections · 10 reactions · 5 readers                       │
│  │  Avg read-through: 60%                                       │
│  │  Changes: Rewrote "Adoption Cliff" to add examples,          │
│  │           shortened "Governance" section                      │
│  │  [View]  [View Feedback]  [Restore]                          │
│  │                                                               │
│  ● v2 ─── Feb 28, 2026 ────────────────────────────────────    │
│  │  5 sections · 7 reactions · 4 readers                        │
│  │  Avg read-through: 48%                                       │
│  │  Changes: Added "Making It Work" section,                     │
│  │           fixed intro based on v1 feedback                    │
│  │  [View]  [View Feedback]  [Restore]                          │
│  │                                                               │
│  ● v1 ─── Feb 12, 2026 ────────────────────────────────────    │
│     4 sections · 4 reactions · 3 readers                        │
│     Avg read-through: 32%                                       │
│     Initial upload from Google Docs                              │
│     [View]  [View Feedback]  [Restore]                          │
│                                                                  │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │ 📈 Trend: Avg read-through improved 32% → 48% → 60%     │  │
│  │    Each revision is getting readers further.              │  │
│  │    Your next draft-killer is "The Adoption Cliff."       │  │
│  └───────────────────────────────────────────────────────────┘  │
│                                                                  │
│                          [Upload New Version]                    │
└──────────────────────────────────────────────────────────────────┘
```

**Functionality:** Timeline of all versions with key metrics per version.
Each version preserves its own feedback, so you can see how reactions
changed across revisions. "Restore" makes an old version current (with
confirmation). "View" opens a read-only snapshot. The trend box at the
bottom auto-summarizes the trajectory — crucial for morale and for
knowing when you're "done" with beta reading and ready for copy editing.
"Upload New Version" triggers the same import flow as New Article but
creates a new version under the same article.

```yaml
screen: version_history
window:
  component: MacWindow
  title: "Version History"
  size: [600, 520]

layout:
  - heading: "{article.title}"

  - timeline:
      style: vertical  # connected dots with lines
      items: article.versions.reverse()  # newest first
      item_template:
        - dot: "● if current else ○"
        - row:
            - label: "v{version}" + "(current)" if current
            - date: "{created_at, formatted}"
        - stats: "{sections} sections · {reactions} reactions · {readers} readers"
        - metric: "Avg read-through: {avg_progress}%"
        - text: "Changes: {changelog}"  # author-written
        - action_row:
            - link: { text: "View", action: open_readonly_snapshot }
            - link: { text: "View Feedback", action: navigate(/article/{id}/feedback?version={v}) }
            - link: { text: "Restore", action: confirm_then(article.restore_version), hidden_if: current }

  - insight_box:
      icon: "📈"
      style: hatched_background
      content:
        template: "Trend: Avg read-through improved {v1}% → {v2}% → {v3}%.\n{interpretation}"
        interpretation_rules:
          improving: "Each revision is getting readers further."
          plateau: "Progress has leveled off — the current draft-killer needs a bigger change."
          declining: "Last revision may have introduced new issues. Check the feedback."

  - button: { label: "Upload New Version", primary: true, action: show_import_modal }
```

---

## 15 · READER MANAGEMENT

```
┌──────────────────────────────────────────────────────────────────┐
│☐│▒▒▒▒▒▒▒▒▒▒▒▒▒ Manage Readers ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒│
├──────────────────────────────────────────────────────────────────┤
│  Search: [________________]  Filter: [All Articles ▾]           │
│──────────────────────────────────────────────────────────────────│
│                                                                  │
│ ┌────┬──────────┬──────────┬────────┬───────────┬────────────┐  │
│ │    │ Name     │ Email    │Articles│ Reactions │ Status     │  │
│ ├────┼──────────┼──────────┼────────┼───────────┼────────────┤  │
│ │ SK │ Sarah K. │ s@e.com  │   2    │    7      │ ● Active   │  │
│ │ MT │ Marcus T.│ m@e.com  │   1    │    3      │ ● Active   │  │
│ │ PR │ Priya R. │ p@e.com  │   1    │    4      │ ● Active   │  │
│ │ JL │ James L. │ j@e.com  │   2    │    2      │ ○ Stale    │  │
│ │ CW │ Chen W.  │ c@e.com  │   1    │    1      │ ◐ Invited  │  │
│ └────┴──────────┴──────────┴────────┴───────────┴────────────┘  │
│                                                                  │
│  ─── Selected: Sarah K. ────────────────────────────────────    │
│                                                                  │
│  ┌ Articles ─────────────────────────────────────────────────┐  │
│  │ Why Design Systems Fail    100% read    5 reactions       │  │
│  │ Remote Work Isn't Working  100% read    1 reaction        │  │
│  └───────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌ Recent Reactions ─────────────────────────────────────────┐  │
│  │ ★ "Great hook" — Intro, Design Systems                   │  │
│  │ ♥ "YES. Exactly what happened" — Adoption, Design Sys    │  │
│  │ ★ "2-minute rule is a great benchmark" — Docs, Design    │  │
│  └───────────────────────────────────────────────────────────┘  │
│                                                                  │
│  [Invite to Another Article]  [Send Follow-up]  [Remove]        │
└──────────────────────────────────────────────────────────────────┘
```

**Functionality:** Cross-article view of all beta readers. Table shows
each reader's engagement across all articles. Click a reader to see a
detail panel with per-article breakdown and recent reactions. Status
is auto-computed: Active (left feedback recently), Stale (invited but
inactive for 7+ days), Invited (link sent, not yet opened). Actions:
invite them to another article, send a follow-up nudge email, or remove
from all articles. Follows the beta reading guide's emphasis on tracking
"engaged" vs "invited" readers and the ~1-in-4 engagement rule.

```yaml
screen: reader_management
window:
  component: MacWindow
  title: "Manage Readers"
  maximized: true

toolbar:
  - search: { placeholder: "Search…", field: name_or_email }
  - dropdown: { id: article_filter, label: "Filter", options: [All Articles, ...articles] }

main:
  component: master_detail  # click row to show detail panel below

  table:
    columns:
      - { id: avatar, width: 40, render: circle_avatar }
      - { id: name, label: "Name", sortable: true }
      - { id: email, label: "Email", sortable: true }
      - { id: articles_count, label: "Articles", sortable: true }
      - { id: reactions_count, label: "Reactions", sortable: true }
      - { id: status, label: "Status", sortable: true, render: status_dot }
    status_rules:
      active: "left feedback within 7 days"
      stale: "invited 7+ days ago, no recent activity"
      invited: "link sent, not yet opened"

  detail_panel:
    visible_when: row_selected
    sections:
      - list:
          label: "Articles"
          items: reader.articles
          template: "{article.title}    {progress}% read    {reactions} reactions"
      - list:
          label: "Recent Reactions"
          items: reader.reactions.recent(5)
          template: "{icon} \"{text}\" — {section}, {article}"
      - action_row:
          - button: { label: "Invite to Another Article", action: show_invite_to_article_picker }
          - button: { label: "Send Follow-up", action: send_followup_email }
          - button: { label: "Remove", style: danger, action: confirm_then(reader.remove_all) }
```

---

## 16 · READER WELCOME / LANDING (public link page)

```
┌──────────────────────────────────────────────────────────────────┐
│ 🍎  Draft Review                                                │
├──────────────────────────────────────────────────────────────────┤
│░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░│
│░░░░  ┌──────────────────────────────────────────────┐  ░░░░░░░░░│
│░░░░  │☐│▒▒▒▒▒▒▒▒▒ You've Been Invited ▒▒▒▒▒▒▒▒▒▒▒▒│  ░░░░░░░░░│
│░░░░  ├──────────────────────────────────────────────┤  ░░░░░░░░░│
│░░░░  │                                              │  ░░░░░░░░░│
│░░░░  │  Alex Chen is looking for your feedback on:  │  ░░░░░░░░░│
│░░░░  │                                              │  ░░░░░░░░░│
│░░░░  │  ┌────────────────────────────────────────┐  │  ░░░░░░░░░│
│░░░░  │  │                                        │  │  ░░░░░░░░░│
│░░░░  │  │     📝  Why Design Systems Fail        │  │  ░░░░░░░░░│
│░░░░  │  │     Draft 3 — March 2026               │  │  ░░░░░░░░░│
│░░░░  │  │     5 sections · ~8 min read            │  │  ░░░░░░░░░│
│░░░░  │  │                                        │  │  ░░░░░░░░░│
│░░░░  │  └────────────────────────────────────────┘  │  ░░░░░░░░░│
│░░░░  │                                              │  ░░░░░░░░░│
│░░░░  │  ┌─ A note from the author ───────────────┐  │  ░░░░░░░░░│
│░░░░  │  │ Thanks for reading an early draft! I'm │  │  ░░░░░░░░░│
│░░░░  │  │ especially looking for feedback on      │  │  ░░░░░░░░░│
│░░░░  │  │ what's useful and where things get      │  │  ░░░░░░░░░│
│░░░░  │  │ confusing.                              │  │  ░░░░░░░░░│
│░░░░  │  └─────────────────────────────────────────┘  │  ░░░░░░░░░│
│░░░░  │                                              │  ░░░░░░░░░│
│░░░░  │  ─── If you have an account: ───────────── │  ░░░░░░░░░░│
│░░░░  │  ╔════════════════════════════════════════╗  │  ░░░░░░░░░│
│░░░░  │  ║         Sign In & Start Reading        ║  │  ░░░░░░░░░│
│░░░░  │  ╚════════════════════════════════════════╝  │  ░░░░░░░░░│
│░░░░  │                                              │  ░░░░░░░░░│
│░░░░  │  ─── New here? ────────────────────────── │  ░░░░░░░░░░│
│░░░░  │  ┌────────────────────────────────────────┐  │  ░░░░░░░░░│
│░░░░  │  │      Create Account & Start Reading     │  │  ░░░░░░░░░│
│░░░░  │  └────────────────────────────────────────┘  │  ░░░░░░░░░│
│░░░░  │                                              │  ░░░░░░░░░│
│░░░░  │   or  [Read as Guest] (limited features)     │  ░░░░░░░░░│
│░░░░  │                                              │  ░░░░░░░░░│
│░░░░  └──────────────────────────────────────────────┘  ░░░░░░░░░│
│░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░│
└──────────────────────────────────────────────────────────────────┘
```

```yaml
screen: reader_landing
window:
  component: MacWindow
  title: "You've Been Invited"
  centered: true
  size: [480, 540]

layout:
  - text: "{author.name} is looking for your feedback on:"

  - article_preview_card:
      border: "2px solid P.black"
      padding: 20
      center: true
      children:
        - icon: { emoji: "📝", size: 36, border: true }
        - title: "{article.title}"
        - subtitle: "Draft {version} — {date}"
        - meta: "{sections} sections · ~{estimated_reading_time} min read"

  - author_note_box:
      label: "A note from the author"
      text: "{article.author_note}"
      style: bordered_inset

  - divider: { text: "If you have an account" }
  - button: { label: "Sign In & Start Reading", primary: true, action: navigate(/login?redirect=/read/{id}) }

  - divider: { text: "New here?" }
  - button: { label: "Create Account & Start Reading", action: navigate(/signup?invite={token}) }

  - conditional:
      if: "article.access_mode == 'link'"
      show:
        link: { text: "or Read as Guest (limited features)", action: start_guest_session }

  - conditional:
      if: "article.access_mode == 'password'"
      show:
        field: { id: password, label: "Enter article password", type: password }
```

---

## 17 · READER VIEW — ✅ ALREADY BUILT

> Implemented in `reader-view.jsx`. See existing code.
> Includes: Welcome splash → section-by-section reading → inline
> paragraph reactions (★?◎♥) → progress dots → completion dialog.

---

## 18 · READER SUMMARY (POST-READ)

```
┌──────────────────────────────────────────────────────────────────┐
│☐│▒▒▒▒▒▒▒▒▒▒▒▒ Review Complete ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒│
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│                          🎉                                      │
│                  Thank you for reading!                           │
│                                                                  │
│  Your Reactions ──────────────────────────────────────────────   │
│                                                                  │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌───────────┐ │
│  │     ★ 3     │ │     ? 1     │ │     ◎ 1     │ │    ♥ 2    │ │
│  │   Useful    │ │  Confused   │ │    Slow     │ │ Favorite  │ │
│  └─────────────┘ └─────────────┘ └─────────────┘ └───────────┘ │
│                                                                  │
│  You read 5 of 5 sections (100%)                                │
│  ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ 100%          │
│                                                                  │
│  ┄┄┄ Optional Follow-up ┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄   │
│                                                                  │
│  Any overall thoughts for the author?                            │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │ Really enjoyed it. The section on governance was the most │  │
│  │ actionable part. Would love more real-world examples      │  │
│  │ throughout.                                                │  │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  Would you recommend this to a colleague?                        │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐                        │
│  │ 😐 Maybe │ │ 👍 Yes   │ │ 🔥 100%! │                        │
│  └──────────┘ └──────────┘ └──────────┘                        │
│                                                                  │
│  [ ] Notify me when the next version is ready                    │
│                                                                  │
│            ╔══════════════════════════╗                          │
│            ║    Submit & Close        ║                          │
│            ╚══════════════════════════╝                          │
│                                                                  │
│  Your feedback helps make this article better.                   │
│  Every reaction matters. Thanks for being a reader.              │
└──────────────────────────────────────────────────────────────────┘
```

**Functionality:** Shown after the reader hits "Finish Review." Displays
summary stats of their reactions. Optional open-ended textarea for
overall thoughts. Recommendability score (the ultimate test from the
beta reading guide — would you tell someone else to read this?).
Opt-in checkbox to receive a notification when a new version is
available. The recommendability question is directly from the
Write Useful Books methodology.

```yaml
screen: reader_summary
window:
  component: MacWindow
  title: "Review Complete"
  maximized: true

layout:
  - icon: { emoji: "🎉", size: 48 }
  - heading: "Thank you for reading!"

  - section: "Your Reactions"
    component: stats_grid
    columns: 4
    items: reaction_types.map(t => { icon: t.icon, count: user_reactions[t.type], label: t.label })

  - row:
      - text: "You read {read_count} of {total_count} sections ({percent}%)"
      - component: HatchBar
        percent: read_percent

  - divider: { text: "Optional Follow-up" }

  - field:
      id: overall_thoughts
      label: "Any overall thoughts for the author?"
      type: textarea
      rows: 3
      optional: true

  - selector:
      id: recommendability
      label: "Would you recommend this to a colleague?"
      style: three_cards
      options:
        - { value: maybe, icon: "😐", label: "Maybe" }
        - { value: yes, icon: "👍", label: "Yes" }
        - { value: absolutely, icon: "🔥", label: "100%!" }

  - checkbox:
      id: notify_new_version
      label: "Notify me when the next version is ready"

  - button: { label: "Submit & Close", primary: true, action: submit_summary_and_close }

  - footer_text: "Your feedback helps make this article better.\nEvery reaction matters."
```

---

## 19 · READER HISTORY (MY REVIEWS)

```
┌──────────────────────────────────────────────────────────────────┐
│☐│▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒ My Reviews ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒│
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │ 📝 Why Design Systems Fail              by Alex Chen      │  │
│  │    Read on Mar 12 · v3 · 5/5 sections                     │  │
│  │    ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ 100%               │  │
│  │    Your reactions: ★3  ?1  ◎1  ♥2                         │  │
│  │    ✉ New version available! (v4, Mar 20)                   │  │
│  │               [Re-read New Version]  [View My Reactions]   │  │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │ 📝 Remote Work Isn't Working             by Alex Chen     │  │
│  │    Read on Mar 5 · v1 · 3/3 sections                      │  │
│  │    ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ 100%               │  │
│  │    Your reactions: ♥1  ?1                                  │  │
│  │                                          [View My Reactions] │  │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │ 📝 Scaling Engineering Teams             by Jordan M.     │  │
│  │    Invited Mar 18 · Not yet started                        │  │
│  │    ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ 0%                │  │
│  │                                              [Start Reading] │  │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  You've reviewed 2 articles and left 8 reactions total.          │
└──────────────────────────────────────────────────────────────────┘
```

```yaml
screen: reader_history
window:
  component: MacWindow
  title: "My Reviews"
  maximized: true

list:
  items: reader.articles
  sort: most_recent_first
  item_template:
    component: review_card
    border: "2px solid P.black"
    layout:
      - row:
          - icon: "📝"
          - title: "{article.title}"
          - byline: "by {author.name}"
      - text: "Read on {date} · v{version} · {read_sections}/{total_sections} sections"
      - component: HatchBar
        percent: progress
      - text: "Your reactions: {reaction_summary}"
      - conditional:
          if: "new_version_available"
          show:
            alert: "✉ New version available! (v{new_version}, {new_version_date})"
      - action_row:
          - button:
              label: "Re-read New Version"
              visible_if: new_version_available
              action: navigate(/read/{article_id})
          - button:
              label: "Start Reading"
              visible_if: "progress == 0"
              primary: true
              action: navigate(/read/{article_id})
          - button:
              label: "View My Reactions"
              action: navigate(/my-reactions/{article_id})

footer:
  text: "You've reviewed {count} articles and left {total_reactions} reactions total."
```

---

## 20 · PROFILE & ACCOUNT

```
┌──────────────────────────────────────────────────────────────────┐
│☐│▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒ Account Settings ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒│
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  PROFILE ─────────────────────────────────────────────────────   │
│                                                                  │
│  ┌────┐  Name                                                    │
│  │ AC │  ┌──────────────────────────────────┐                    │
│  └────┘  │ Alex Chen                        │                    │
│          └──────────────────────────────────┘                    │
│  Display   Email                                                 │
│  initials  ┌──────────────────────────────────┐                  │
│  shown to  │ alex@example.com                 │  ✓ Verified      │
│  readers   └──────────────────────────────────┘                  │
│                                                                  │
│  Bio (shown to readers)                                          │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │ Writer & design systems consultant. Building better       │  │
│  │ tools for better teams.                                    │  │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  SECURITY ────────────────────────────────────────────────────   │
│                                                                  │
│  Password             [Change Password]                          │
│  Two-factor auth      [Enable 2FA]                               │
│  Connected accounts   Google ✓ connected  [Disconnect]           │
│                                                                  │
│  PREFERENCES ─────────────────────────────────────────────────   │
│                                                                  │
│  Default author note for new articles:                           │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │ Thanks for reading! Focus on what's useful, confusing, or │  │
│  │ slow. Don't worry about typos.                             │  │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  [✓] Email me when a reader finishes an article                  │
│  [✓] Email me when I get a new reaction                          │
│  [ ] Weekly digest instead of individual notifications            │
│                                                                  │
│  DANGER ZONE ─────────────────────────────────────────────────   │
│  [Delete Account]  This will permanently remove all your data.   │
│                                                                  │
│                           ╔════════════════╗                     │
│                           ║ Save Changes   ║                     │
│                           ╚════════════════╝                     │
└──────────────────────────────────────────────────────────────────┘
```

```yaml
screen: account_settings
window:
  component: MacWindow
  title: "Account Settings"
  size: [560, 620]
  scrollable: true

layout:
  - section: "PROFILE"
    children:
      - row:
          - avatar:
              component: circle_avatar
              initials: "{user.initials}"
              size: 48
              note: "Display initials shown to readers"
          - column:
              - field: { id: name, label: "Name" }
              - field: { id: email, label: "Email", readonly: true, suffix: "✓ Verified" }
      - field: { id: bio, label: "Bio (shown to readers)", type: textarea, rows: 2 }

  - section: "SECURITY"
    children:
      - row: { label: "Password", action_label: "Change Password", action: show_password_modal }
      - row: { label: "Two-factor auth", action_label: "Enable 2FA", action: show_2fa_setup }
      - row:
          label: "Connected accounts"
          content: "Google ✓ connected"
          action_label: "Disconnect"
          action: auth.disconnect_google

  - section: "PREFERENCES"
    children:
      - field:
          id: default_author_note
          label: "Default author note for new articles"
          type: textarea
          rows: 2
      - checkbox_group:
          items:
            - { id: email_finish, label: "Email me when a reader finishes an article", default: true }
            - { id: email_reaction, label: "Email me when I get a new reaction", default: true }
            - { id: weekly_digest, label: "Weekly digest instead of individual notifications", default: false }

  - section: "DANGER ZONE"
    style: danger_border
    children:
      - button:
          label: "Delete Account"
          style: danger
          action: confirm_twice_then(auth.delete_account)
          warning: "This will permanently remove all your data."

  - button: { label: "Save Changes", primary: true }
```

---

## 21 · NOTIFICATION PREFERENCES (integrated in Profile above)

> Notification settings are part of Screen 20 (Profile & Account)
> under the PREFERENCES section. No separate screen needed.
> Email notification types: new reaction, reader finished,
> reader started, weekly digest, new version available (for readers).

---

## 22 · BILLING & SUBSCRIPTION

```
┌──────────────────────────────────────────────────────────────────┐
│☐│▒▒▒▒▒▒▒▒▒▒▒▒▒ Subscription ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒│
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  CURRENT PLAN ────────────────────────────────────────────────   │
│                                                                  │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │ ★ AUTHOR PRO                               $12/month      │  │
│  │   Unlimited articles · Unlimited readers · Analytics       │  │
│  │   Export · Custom reactions · Version history              │  │
│  │   Next billing: April 1, 2026                             │  │
│  │                              [Manage Payment] [Cancel]    │  │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┄┄┄ or compare plans ┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄   │
│                                                                  │
│  ┌─────────────┐  ┌─────────────┐  ┌────────────────────────┐  │
│  │ FREE        │  │ ★ PRO      │  │ ★★ TEAM               │  │
│  │             │  │             │  │                        │  │
│  │ 1 article   │  │ Unlimited   │  │ Everything in Pro      │  │
│  │ 3 readers   │  │ articles    │  │ + 5 author seats       │  │
│  │ Basic       │  │ Unlimited   │  │ + Team analytics       │  │
│  │ reactions   │  │ readers     │  │ + Shared reader pool   │  │
│  │             │  │ Analytics   │  │ + Admin controls       │  │
│  │             │  │ Export      │  │                        │  │
│  │ $0/mo       │  │ $12/mo     │  │ $39/mo                │  │
│  │ [Current]   │  │ [Upgrade]  │  │ [Upgrade]             │  │
│  └─────────────┘  └─────────────┘  └────────────────────────┘  │
│                                                                  │
│  BILLING HISTORY ─────────────────────────────────────────────   │
│                                                                  │
│  Mar 1, 2026   Author Pro   $12.00   ✓ Paid    [Receipt]       │
│  Feb 1, 2026   Author Pro   $12.00   ✓ Paid    [Receipt]       │
│  Jan 1, 2026   Author Pro   $12.00   ✓ Paid    [Receipt]       │
└──────────────────────────────────────────────────────────────────┘
```

```yaml
screen: billing
window:
  component: MacWindow
  title: "Subscription"
  size: [600, 560]

layout:
  - section: "CURRENT PLAN"
    component: plan_card
    highlight: true  # border with double line
    fields:
      - plan_name: "{plan.name}"
      - price: "${plan.price}/month"
      - features: "{plan.features.join(' · ')}"
      - next_billing: "{next_billing_date}"
    actions:
      - link: "Manage Payment"  # opens Stripe portal
      - link: "Cancel"  # confirm dialog with retention offer

  - divider: { text: "or compare plans" }

  - plan_comparison:
      columns: 3
      style: card_grid
      plans:
        - name: "FREE"
          price: "$0/mo"
          features: [1 article, 3 readers, Basic reactions]
          cta: { label: "Current" or "Downgrade", action: plan.change(free) }
        - name: "PRO"
          badge: "★"
          price: "$12/mo"
          features: [Unlimited articles, Unlimited readers, Analytics, Export]
          cta: { label: "Upgrade", primary: true, action: plan.change(pro) }
        - name: "TEAM"
          badge: "★★"
          price: "$39/mo"
          features: [Everything in Pro, 5 author seats, Team analytics, Shared reader pool, Admin controls]
          cta: { label: "Upgrade", action: plan.change(team) }

  - section: "BILLING HISTORY"
    component: table
    columns: [Date, Plan, Amount, Status, Action]
    actions_per_row: [{ label: "Receipt", action: download_receipt }]
```

---

## 23 · TEAM MANAGEMENT (Team plan only)

```
┌──────────────────────────────────────────────────────────────────┐
│☐│▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒ Team Settings ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒│
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  MEMBERS ─────────────────────────────────────── 3 of 5 seats   │
│                                                                  │
│ ┌────┬──────────────┬────────────────┬────────┬──────────────┐  │
│ │    │ Name         │ Email          │ Role   │ Articles     │  │
│ ├────┼──────────────┼────────────────┼────────┼──────────────┤  │
│ │ AC │ Alex Chen    │ alex@team.com  │ Owner  │ 3            │  │
│ │ JM │ Jordan M.    │ jordan@team.co │ Author │ 2            │  │
│ │ LK │ Lisa K.      │ lisa@team.com  │ Author │ 1            │  │
│ └────┴──────────────┴────────────────┴────────┴──────────────┘  │
│                                                                  │
│  [+ Invite Team Member]                                          │
│                                                                  │
│  SHARED READER POOL ──────────────────────────────────────────   │
│                                                                  │
│  12 readers across all team articles                             │
│  [✓] Share reader contacts across team members                   │
│  [ ] Require approval before inviting shared readers             │
│                                                                  │
│  TEAM DEFAULTS ───────────────────────────────────────────────   │
│                                                                  │
│  Default access mode:  [Invite-only ▾]                           │
│  Default reactions:    [✓]★  [✓]?  [✓]◎  [✓]♥                   │
│  [✓] New members inherit team defaults                           │
│                                                                  │
│                           ╔════════════════╗                     │
│                           ║ Save Settings  ║                     │
│                           ╚════════════════╝                     │
└──────────────────────────────────────────────────────────────────┘
```

```yaml
screen: team_management
window:
  component: MacWindow
  title: "Team Settings"
  size: [580, 520]
  requires: team_plan

layout:
  - section: "MEMBERS"
    header_suffix: "{member_count} of {seat_limit} seats"
    component: table
    columns:
      - { id: avatar, width: 40 }
      - { id: name, label: "Name" }
      - { id: email, label: "Email" }
      - { id: role, label: "Role", editable: dropdown, options: [Owner, Author, Viewer] }
      - { id: articles, label: "Articles" }
    row_actions: [Remove]  # except Owner
    footer:
      button: { label: "+ Invite Team Member", action: show_team_invite_modal }

  - section: "SHARED READER POOL"
    children:
      - text: "{total_readers} readers across all team articles"
      - checkbox: { id: share_readers, label: "Share reader contacts across team members" }
      - checkbox: { id: require_approval, label: "Require approval before inviting shared readers" }

  - section: "TEAM DEFAULTS"
    children:
      - dropdown: { id: default_access, label: "Default access mode", options: [Invite-only, Anyone with link, Password] }
      - reaction_toggles: { label: "Default reactions", items: all_reaction_types }
      - checkbox: { id: inherit_defaults, label: "New members inherit team defaults" }

  - button: { label: "Save Settings", primary: true }
```

---

## 24 · ADMIN OVERVIEW (Team plan only)

```
┌──────────────────────────────────────────────────────────────────┐
│☐│▒▒▒▒▒▒▒▒▒▒▒▒▒▒ Team Overview ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒│
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌──────────┐ │
│  │ 6       │ │ 43      │ │ 12      │ │ 58%     │ │ 3        │ │
│  │ARTICLES │ │REACTIONS │ │READERS  │ │AVG READ │ │AUTHORS   │ │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘ └──────────┘ │
│                                                                  │
│  ARTICLES ACROSS TEAM ────────────────────────────────────────   │
│                                                                  │
│  Author     Article                   Readers  React  Avg Read  │
│  ─────────────────────────────────────────────────────────────   │
│  Alex C.    Why Design Systems Fail     5       10     60%      │
│  Alex C.    Remote Work Isn't Working   2        2     66%      │
│  Alex C.    The Art of Saying No        0        0      —       │
│  Jordan M.  Scaling Engineering Teams   4        8     45%      │
│  Jordan M.  On-call Best Practices      3        5     72%      │
│  Lisa K.    Writing Better RFCs         2        3     50%      │
│                                                                  │
│  TEAM ACTIVITY ───────────────────────────────────────────────   │
│                                                                  │
│  This week: +12 reactions, +3 new readers, 1 new article        │
│  Last week: +8 reactions, +2 new readers, 0 new articles        │
│                                                                  │
│            [Export Team Report]  [Team Settings →]               │
└──────────────────────────────────────────────────────────────────┘
```

```yaml
screen: admin_overview
window:
  component: MacWindow
  title: "Team Overview"
  maximized: true
  requires: team_plan + owner_role

layout:
  - stats_row:
      items:
        - { label: "ARTICLES", value: team.total_articles }
        - { label: "REACTIONS", value: team.total_reactions }
        - { label: "READERS", value: team.total_readers }
        - { label: "AVG READ", value: team.avg_read_percent + "%" }
        - { label: "AUTHORS", value: team.member_count }

  - section: "ARTICLES ACROSS TEAM"
    component: table
    sortable: true
    columns:
      - { id: author, label: "Author" }
      - { id: title, label: "Article" }
      - { id: readers, label: "Readers", sortable: true }
      - { id: reactions, label: "React", sortable: true }
      - { id: avg_read, label: "Avg Read", sortable: true }
    on_click_row: navigate(/article/{id}/analytics)

  - section: "TEAM ACTIVITY"
    component: week_over_week
    metrics: [reactions, new_readers, new_articles]
    format: "This week: {metrics} / Last week: {metrics}"

  - action_row:
      - button: { label: "Export Team Report", action: export.team_report }
      - button: { label: "Team Settings →", action: navigate(/team-settings) }
```

---

## 25 · 404 / ERROR

```
┌─────────────────────────────────────────┐
│☐│▒▒▒▒▒▒▒▒▒▒▒ Error ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒│
├─────────────────────────────────────────┤
│                                         │
│                  💀                      │
│                                         │
│          Sorry, System Error            │
│                                         │
│   The application has encountered an    │
│   unexpected error. The document you    │
│   requested could not be found.         │
│                                         │
│   Error: 404 — Not Found               │
│                                         │
│          ID: 10093F87-AB21              │
│                                         │
│  ╔═══════════════════════════════════╗  │
│  ║          Return Home              ║  │
│  ╚═══════════════════════════════════╝  │
│                                         │
│          [Report This Problem]          │
└─────────────────────────────────────────┘
```

```yaml
screen: error
window:
  component: MacWindow
  title: "Error"
  centered: true
  size: [380, 340]

variants:
  404:
    icon: "💀"
    heading: "Sorry, System Error"
    message: "The document you requested could not be found."
    code: "Error: 404 — Not Found"
  403:
    icon: "🔒"
    heading: "Access Denied"
    message: "You don't have permission to view this article. Ask the author for an invitation."
  500:
    icon: "💣"
    heading: "Something Went Wrong"
    message: "We hit an unexpected error. We've been notified and are looking into it."
  offline:
    icon: "📡"
    heading: "No Connection"
    message: "Check your internet connection and try again."

layout:
  - icon: { emoji: variant.icon, size: 40 }
  - heading: "{variant.heading}"
  - text: "{variant.message}"
  - code: "{variant.code}"
  - text: "ID: {error_id}" style: monospace, small
  - button: { label: "Return Home", primary: true, action: navigate(/) }
  - link: { text: "Report This Problem", action: open_support_form }
```

---

## 26 · ONBOARDING TOUR (first-time author)

```
┌──────────────────────────────────────────────────────────────────┐
│☐│▒▒▒▒▒▒▒▒▒▒▒▒▒ Welcome to Draft Review ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒│
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│     ●───────────○───────────○───────────○                        │
│   Upload      Invite      Collect     Revise                    │
│                                                                  │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │                                                            │  │
│  │                        📄                                  │  │
│  │                                                            │  │
│  │              Step 1: Upload Your Draft                     │  │
│  │                                                            │  │
│  │  Start by uploading your article, essay, or chapter.       │  │
│  │  We support Word docs, Google Docs, and Markdown.          │  │
│  │                                                            │  │
│  │  We'll split it into sections by heading so your           │  │
│  │  readers can react to each part individually.              │  │
│  │                                                            │  │
│  │  Don't wait until it's perfect — early, messy drafts       │  │
│  │  are where beta reading has the most leverage.             │  │
│  │                                                            │  │
│  │  ┌──────────────────────────────────────────────────────┐  │  │
│  │  │  💡 Tip: You only need 3-4 engaged readers to       │  │  │
│  │  │  direct your next revision effectively.              │  │  │
│  │  └──────────────────────────────────────────────────────┘  │  │
│  │                                                            │  │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│                                                                  │
│     [Skip Tour]                          [Next: Invite →]       │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
```

**Functionality:** Four-step onboarding wizard shown once after first
sign-up. Steps: (1) Upload your draft, (2) Invite 3-4 readers,
(3) Collect reactions, (4) Revise and repeat. Each step has an
illustration, explanation, and a beta-reading tip from the methodology.
Can be skipped. After completion, lands on the New Article screen.
Progress dots at top show current step.

```yaml
screen: onboarding
window:
  component: MacWindow
  title: "Welcome to Draft Review"
  maximized: true

progress:
  component: step_dots
  steps: [Upload, Invite, Collect, Revise]
  current: 0

steps:
  - id: upload
    icon: "📄"
    heading: "Step 1: Upload Your Draft"
    body: "Start by uploading your article, essay, or chapter. We support Word docs, Google Docs, and Markdown.\n\nWe'll split it into sections by heading so your readers can react to each part individually.\n\nDon't wait until it's perfect — early, messy drafts are where beta reading has the most leverage."
    tip: "You only need 3-4 engaged readers to direct your next revision effectively."

  - id: invite
    icon: "✉"
    heading: "Step 2: Invite Readers"
    body: "Think of one person who would genuinely want to read what you're writing. Send them a personal invitation.\n\nDo this once a day for two weeks, and you'll have 3-4 engaged readers."
    tip: "Only ~1 in 4 invitees will engage. That's normal. Don't take it personally."

  - id: collect
    icon: "✦"
    heading: "Step 3: Collect Reactions"
    body: "Your readers will mark what's useful, confusing, or slow as they read. You'll see where they stop, what they love, and where they get lost.\n\nLook for the draft-killer: the section causing the most drop-off."
    tip: "Readers spot the problems. You design the solutions. Ignore their specific fix suggestions."

  - id: revise
    icon: "✎"
    heading: "Step 4: Revise & Repeat"
    body: "Use the feedback to revise your draft, focusing on the biggest issue first. Then upload the new version and invite more readers.\n\nEach round, your article gets better and readers get further."
    tip: "Once most readers reach the end and get value, you're ready for copy editing."

navigation:
  - link: { text: "Skip Tour", action: navigate(/dashboard), style: subtle }
  - spacer: flex
  - button: { label: "Next: {next_step} →", primary: true }
  # Last step button: "Create My First Article →" → navigate(/new-article)
```

---

## 27 · KEYBOARD SHORTCUTS (MODAL)

```
┌─────────────────────────────────────────┐
│☐│▒▒▒▒▒▒▒ Keyboard Shortcuts ▒▒▒▒▒▒▒▒▒▒│
├─────────────────────────────────────────┤
│                                         │
│  NAVIGATION                             │
│  ← →        Previous / Next section     │
│  ⌘ 1-9      Jump to section by number   │
│  ⌘ ↑        Back to Dashboard           │
│  Tab        Next paragraph              │
│                                         │
│  REACTIONS (while paragraph focused)    │
│  1          ★ Useful                    │
│  2          ? Confused                  │
│  3          ◎ Slow                      │
│  4          ♥ Favorite                  │
│  Enter      Submit reaction             │
│  Esc        Cancel                      │
│                                         │
│  GENERAL                                │
│  ⌘ N        New article                 │
│  ⌘ E        Export feedback             │
│  ⌘ I        Invite reader              │
│  ⌘ /        This help                  │
│  ⌘ ,        Settings                   │
│                                         │
│          ╔══════════════════╗           │
│          ║       OK         ║           │
│          ╚══════════════════╝           │
└─────────────────────────────────────────┘
```

```yaml
screen: keyboard_shortcuts
type: modal
component: MacWindow
title: "Keyboard Shortcuts"
size: [380, 420]

layout:
  - shortcut_group:
      label: "NAVIGATION"
      items:
        - { keys: "← →", desc: "Previous / Next section" }
        - { keys: "⌘ 1-9", desc: "Jump to section by number" }
        - { keys: "⌘ ↑", desc: "Back to Dashboard" }
        - { keys: "Tab", desc: "Next paragraph" }

  - shortcut_group:
      label: "REACTIONS (while paragraph focused)"
      items:
        - { keys: "1", desc: "★ Useful" }
        - { keys: "2", desc: "? Confused" }
        - { keys: "3", desc: "◎ Slow" }
        - { keys: "4", desc: "♥ Favorite" }
        - { keys: "Enter", desc: "Submit reaction" }
        - { keys: "Esc", desc: "Cancel" }

  - shortcut_group:
      label: "GENERAL"
      items:
        - { keys: "⌘ N", desc: "New article" }
        - { keys: "⌘ E", desc: "Export feedback" }
        - { keys: "⌘ I", desc: "Invite reader" }
        - { keys: "⌘ /", desc: "This help" }
        - { keys: "⌘ ,", desc: "Settings" }

  - button: { label: "OK", primary: true }
```

---

## COMPONENT REFERENCE

Summary of all reusable components from the existing codebase,
plus new ones defined in this spec:

```yaml
existing_components:
  MacWindow:      { props: [title, x, y, w, h, onClose, zIndex, maximized] }
  TitleBar:       { props: [title, onClose, extra] }
  Btn:            { props: [children, onClick, primary, small, disabled, style] }
  MenuBar:        { props: [onNewArticle, onAbout, currentView, setView] }
  MenuDropdown:   { props: [items, onClose] }
  ProgressBar:    { props: [percent, width, height] }
  HatchBar:       { props: [percent, w, h] }
  SectionNav:     { props: [sections, current, readSections, onPick] }
  Paragraph:      { props: [text, pid, reactions, onReact, onRemove] }
  ReactionBadge:  { props: [type, text, onRemove] }
  InviteDialog:   { props: [onClose] }
  AboutDialog:    { props: [onClose] }
  DoneDialog:     { props: [stats, onClose] }
  WelcomeSplash:  { props: [article, onStart] }

new_components_needed:
  mac_input:         { desc: "Styled input: 2px black border, GENEVA font, white bg" }
  mac_dropdown:      { desc: "Styled select: 2px border, CHICAGO font, inverted on open" }
  mac_textarea:      { desc: "Styled textarea: same as mac_input" }
  circle_avatar:     { desc: "24-48px circle with initials, 1.5px black border" }
  status_pill:       { desc: "Small rounded badge: [Draft] [In Review] [New] etc" }
  card:              { desc: "2px black border container with optional header" }
  sortable_list:     { desc: "Drag-reorderable list with ≡ handles" }
  table:             { desc: "Retro table: 1.5px borders, CHICAGO headers, dotted row dividers" }
  code_input:        { desc: "6-digit verification boxes" }
  step_dots:         { desc: "Horizontal ●○○○ progress indicator" }
  dropoff_chart:     { desc: "Horizontal bar chart showing reader reach per section" }
  heatmap_grid:      { desc: "Grid of blocks: rows=sections, cols=reaction types" }
  insight_box:       { desc: "Hatched-background callout box with icon" }
  shortcut_row:      { desc: "Key combo + description, monospace keys" }
  plan_card:         { desc: "Pricing card with features list and CTA" }
  feedback_card:     { desc: "Reaction card with type icon, reader, text, status dropdown" }
  week_over_week:    { desc: "Two-line metric comparison (this week vs last)" }
  master_detail:     { desc: "Table above, detail panel below on selection" }
  file_drop_zone:    { desc: "Dashed border drop area for file upload" }

shared_styles:
  P:
    bg: "#e8e8e8"
    white: "#ffffff"
    black: "#000000"
    gray: "#a0a0a0"
    dark: "#555555"
    light: "#c0c0c0"
    stripes: "repeating-linear-gradient(0deg,#fff 0px,#fff 1px,#000 1px,#000 2px)"
  fonts:
    CHICAGO: '"Chicago_12","ChicagoFLF","Geneva",monospace'
    GENEVA: '"Geneva","Monaco",monospace'
```

---

## NAVIGATION MAP

```
                    ┌──────────┐
                    │  Login   │◄──── Forgot Password
                    │  Signup  │◄──── Email Verify
                    └────┬─────┘
                         │
              ┌──────────┴──────────┐
              │                     │
        ┌─────▼─────┐        ┌─────▼─────┐
        │  AUTHOR   │        │  READER   │
        │ Dashboard │        │  History  │
        └─────┬─────┘        └─────┬─────┘
              │                     │
    ┌─────────┼─────────┐     ┌────▼─────┐
    │         │         │     │  Reader  │
┌───▼──┐ ┌───▼──┐ ┌────▼──┐  │  View   │
│Article│ │Reader│ │Account│  │(reading)│
│Manager│ │Mgmt  │ │Billing│  └────┬────┘
└───┬───┘ └──────┘ └───────┘       │
    │                          ┌───▼────┐
┌───▼────────┐                 │ Reader │
│ New Article│                 │Summary │
│ Editor     │                 └────────┘
│ Settings   │
│ Analytics  │
│ Feedback   │
│ Versions   │
│ Export     │
└────────────┘
```

---

*Total screens: 27 (3 already built, 24 specified)*
*New components needed: 18*
*This spec covers the complete product from auth → onboarding → authoring → reading → analytics → team management → billing.*
