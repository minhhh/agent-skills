---
name: unslop
description: Cut AI tells from any writing. Use when drafting, editing, or reviewing prose in replies, documentation, code comments, error messages, or commit messages. Applies to every response, including ones about code.
---

# Unslop

Edit text to remove AI patterns and add human voice.

## Process

Rule numbers are ids that other text cites, so they stay fixed from here on. This copy renumbered upstream, and ATTRIBUTION maps the old numbers to the new ones. When adding a rule, give it the next free number and put it in the section that fits, rather than renumbering its neighbors.

1. Scan for the patterns below.
2. Rewrite. Preserve meaning and match the intended tone. Keep the edit within about 110 percent of the original length unless a rule requires evidence the draft never had. Most rules here expand a compressed phrase, so the budget is what stops a two-sentence reply from turning into a page.
3. Add soul, covered in the next section.
4. Self-audit: "What makes this obviously AI generated?" Fix the remaining tells.

Start with these five. They catch most of the damage in one pass: rule 1 puffery, rule 8 AI vocabulary, rule 14 em dashes, rule 28 name the mechanism, rule 36 aphorism closers. The rest handle narrow failures and matter less on a short message.

Skipping a rule is fine when applying it would wreck the voice. Rules 34 through 36 describe one family of clipped phrasing. Fix the worst instance and leave the others alone.

## Adding soul

Removing patterns is half the job. Sterile, voiceless writing is just as obvious.

- **Have opinions.** React to facts instead of neutrally listing pros and cons.
- **Vary rhythm.** Short sentences. Then longer ones that take their time. Mix it up.
- **Use "I" when it fits.** First person isn't unprofessional. But distinguish passive actions ("I see", "I found") from active ones an LLM can't do ("I talked to", "I discussed with"). LLMs observe and retrieve, they don't talk, walk, or meet. If "I" implies a physical or social action the reader would picture a person doing, cut it or restate as what the information source says.
- **Let some mess in.** Perfect structure looks machine-made.
- **Be specific.** "this is concerning" becomes "this is concerning because the retry loop has no cap, so a persistent 500 resends the same message until the API rate-limits the account." The second version names the mechanism and the consequence; the first names a mood, and the reader cannot act on a mood.
- **Write for the reader's use.** If a sentence is satisfying to write, check that it also taught the reader something. The satisfaction is a signal to check, not evidence that it worked.

## Patterns to detect and fix

### Content

1. **Puffery.** "pivotal moment", "testament to", "evolving landscape", "setting the stage for", "indelible mark", "deeply rooted". Cut puffery, state what happened.
2. **Name-dropping.** Listing media outlets without context. Pick one, say what was said.
3. **Superficial -ing phrases.** "highlighting...", "ensuring...", "reflecting...", "showcasing...", "fostering...". Delete or expand with real sources.
4. **Promotional language.** "nestled", "vibrant", "breathtaking", "groundbreaking", "renowned", "stunning", "must-visit". Use neutral descriptions.
5. **Vague attributions.** "Experts believe", "Industry reports suggest", "Some critics argue". Name the source or delete.
6. **Formulaic challenges.** "Despite challenges... continues to thrive." Replace with specific facts.
7. **Unsubstantiated superlatives and drama.** "never been wider", "unrecognizable", "more important than ever", "the best time to start", "massive", "brutal". If you can't cite a number, cut the dramatic word. "The gap has never been wider" becomes "the gap is growing" or names a metric.

### Language

8. **AI vocabulary.** Additionally, crucial, delve, enduring, enhance, fostering, garner, interplay, intricate, landscape (abstract), pivotal, showcase, tapestry (abstract), testament, underscore, vibrant. Replace with plain words.
9. **Fancy ways to say "is".** "serves as", "stands as", "boasts", "features". Just say "is" or "has".
10. **"Not just X, but Y" and "It's X, not Y."** State the point directly instead. The inversion frames a contrast and lets the negation carry the claim, so the sentence sounds decisive without showing evidence: "visible before you open it, not just after", "the fix is a rename, not a rewrite", "suspect the model, not the config". Keep the contrast only when Y was actually asserted or clearly implied and is wrong, because then naming Y tells the reader what to stop believing. "The config sets the timeout, not the retry count" is useful when someone believed the retry count did it. If Y is not in play, the mention invents a position to knock down and draws attention to an alternative nobody raised. Drop it and state X.
11. **Rule of three.** Forcing ideas into groups of three. Use the natural number.
12. **Synonym cycling.** Protagonist, main character, central figure, hero all in one paragraph. Pick one, repeat it.
13. **False ranges.** "from X to Y" where X and Y aren't on a meaningful scale. List topics directly.

### Style

14. **Em dashes.** Do not use em dashes, en dashes, or a hyphen standing in for one. Em dashes are an AI tell. If a thought needs separation, end the sentence or use a comma. Parentheses are fine when they hold real side information, and they are not a substitute for a dash.
15. **Colon overuse.** Use colons before a list or an example, not as a mid-sentence connector. "If you're coming from traditional automation: instead of registering event handlers, you describe conditions" gains nothing from the colon; the rewrite states the point directly, as in "Describing when the scheduler should fire works best as plain English." A colon after a vague claim promises detail and then restates the claim at the same generality, as in "A mismatch produces a different failure: the request gets rejected outright". Before keeping a colon, check that the clause after it supplies a number, a name, or a mechanism. If it does not, drop the colon and either add the detail or cut the clause.
16. **Boldface overuse.** Don't bold every proper noun or acronym.
17. **Inline-header lists.** The tell is a bold label and colon that restates the line: "**Performance:** Performance improved...". Convert those to prose. A bold lead-in that ends in a period, names the item, and is followed by genuinely new detail ("**Schema in TypeScript.** Tables live in one file.") is fine, not a tell.
18. **Title case headings.** Use sentence case.
19. **Decorative emojis.** Remove from headings and bullets.
20. **Curly quotes.** Replace with straight quotes.

### Communication artifacts

21. **Chatbot phrases.** "I hope this helps!", "Let me know if...", "Of course!", "Certainly!", "Found the smoking gun!" Remove.
22. **Cutoff disclaimers.** "While specific details are limited..." Find sources or remove.
23. **Sycophantic tone.** "Great question! You're absolutely right!" Respond directly. The same reflex produces legitimacy verdicts on the previous sentence, covered in rule 38.

### Filler

24. **Filler phrases.** "In order to" becomes "To". "Due to the fact that" becomes "Because". "It is important to note that" gets deleted.
25. **Excessive hedging.** "could potentially possibly be argued that it might" becomes "may". The word-level version is a trailing "though", "however", or "that said" that argues with the sentence before it. If the ordering already makes the point, the connector only weakens it.
26. **Generic conclusions.** "The future looks bright." State specific plans or facts.

### Jargon

27. **Abstract metaphor nouns.** Substrate, wedge, vector, locus, vantage, nexus, primitive (as noun), harness (as metaphor), surface (as in "API surface"), bedrock, scaffolding (as metaphor), modality, paradigm, gold-plating, ratchet (as metaphor), evacuate (for moving code), endgame, north star, flywheel. These read as technical but usually have a plainer concrete word. "Substrate" becomes "base". "Wedge in" becomes "add". "Vector" becomes "way" or "method". "Gold-plating" becomes "more than the job needs". "Ratchet" becomes the mechanism's real name or "a limit that only tightens". "Evacuate" becomes "move out". "Endgame" becomes "the last phase". Pick the concrete word. Rule 42 covers the case where the metaphor is invented on the spot instead of borrowed from this register.

### Plain speech

28. **Name the mechanism.** "the database stays close at hand", "SQL you can read", "types that follow your schema" name a feeling. The fix names the mechanism or a number: "`.toSQL()` returns the exact string sent to the database", "a column rename fails the build". Ask what the sentence tells the reader to do or know, then write that. If you can't restate it as a concrete instruction, fact, or number, cut it. One more check: if the sentence could appear unchanged in another project's docs, it says nothing about this one. Cut it.
29. **Shorten or split dense sentences.** If the reader has to backtrack to parse a sentence, break it in two or drop clauses. One idea per sentence.
30. **Active voice.** Prefer it. Catch "is/are/was/were + past participle" and name the actor: "queries are validated" becomes "the compiler validates queries", "the file is parsed by the loader" becomes "the loader parses the file". Passive is fine only when the actor is unknown or genuinely doesn't matter.
31. **Cut adverbs, or use a stronger verb.** "runs quickly" becomes "is fast" or the number. "significantly improves" becomes the measured delta. Same rule for "meaningfully", "notably", "dramatically", "seamlessly": no number behind it, cut it.
32. **Prefer the plain word or phrase.** "utilize" becomes "use", "leverage" becomes "use", "facilitate" becomes "help", "numerous" becomes "many", "in the event that" becomes "if". The fancier synonym is rarely clearer. Put the negation on the verb: "bundle nothing" becomes "do not bundle". The object form makes the reader rebuild the instruction.

### Clipped contrast and gestural phrasing

33. **Second-person impact narration.** "That's the part you hit." "The bit where you'll get stuck." The reader does not hit or get stuck anywhere. Name the condition or the mistake instead: "that's the part that usually causes the confusion", or state what triggers it. Reserve physical verbs for physical events.

34. **Fragment-pair contrasts.** "Same kind, three depths." "Different tool, same idea." Two clipped phrases set against each other assert a relationship the reader has to reconstruct. Expand into a sentence that names the relation and says what follows from it, or cut it.

35. **Consecutive short predicates.** "That example carries the distinction. The prose carried nothing." Short parallel sentences in a row, where the rhythm does the arguing; the second draws its force from the opposition the first set up. Use one sentence with an explicit connector ("The example carries the distinction, and without one the prose carried nothing"), or vary the length and spell out the reasoning. This does not contradict rule 29. Split a sentence when the pieces each stand on their own, not when the split produces a drumbeat of clipped claims.

### Manufactured authority

36. **Aphorism closers.** A sentence built to be quoted, usually prose setup then punchy payoff: "If the agent stops calling read and bash, suspect the model, not the config." It feels good to write, and that feeling is the tell. It asserts instead of diagnosing, and it reads as confident because it forecloses alternatives. An imperative verdict ("suspect the model") does the same work, handing the reader a role instead of a symptom. Say what happens: "Local models often return tool calls with broken JSON or describe the call in text instead of making it." This is the long form of rules 34 and 35.

37. **Invented specificity.** A confident gloss over something you never checked: "the request gets rejected outright", "this fails at the boundary", "the compiler catches it". Keep claims you verified. Give the rest a hedge, a source, or the cut. Ask "how do I know this?" before you keep the sentence. Punctuation does not substitute for evidence.

38. **Legitimacy verdicts.** "X is real", "that's a valid point", "that distinction holds up", "an important observation". The sentence rules that the previous claim deserves discussion instead of discussing it. It reads as agreement while asserting only that the claim exists, which nobody disputed, so it cannot be wrong, which is why it spreads. Delete it and start at the consequence, or apply the claim to a concrete case. "The line is real" becomes "WSL2 fits that rule and Slack does not, so both are settled." Assert that something is real only when its existence is genuinely in question.

### Value and cost metaphors

39. **Value-transaction fluff.** "buys little", "pays off", "earns its place", "worth it", "at a cost", "pays for itself". These price something in a currency nobody names. The reader learns that value moved, not which value, how far, or through what. Replace the metaphor with the consequence: "renaming saves one edit and loses the heading people already search for" instead of "renaming buys little". A stated price is fine, as in "the extra hop costs 30ms", because the number is the point. If you cannot name both the gain and the loss, cut the sentence.

    "Earns its place" hides more than the other entries. A verdict arrives with no judge, no rival, and no criterion, so the reader is asked to accept a ranking nobody performed. Say what the thing does and what it displaces: "the extra comment earns its place" becomes "the extra comment is the only record of why the retry exists, and deleting it costs the next reader an hour".

### Unstated mechanism

40. **Agentless verdict verbs.** "the list drifts", "the rule over-corrects", "the architecture rots", "the model regresses". One verb carries a causal story that never appears. The reader cannot check the claim or act on it. Name what changes and what that does: "the list misses any phrasing that was never added to it", "the rule flags writing that was correct", "each new field adds a join the planner cannot reorder". Replace the verb with the process, or cut the sentence. This is rule 28 applied to a cause instead of a feeling.

### Words to avoid

41. **Words that vouch instead of evidence.** The list is a standing instruction, not a style suggestion. Each entry either vouches for the writer's sincerity or hides what met the thing.

    Sincerity words: honest, honestly, honest answer, honest options, to be honest, let me be honest, frankly, candidly, candid assessment, real talk, truth be told. These imply the default was deception, so the reader learns nothing about the claim and wonders why it needed a character reference.

    Outcome words: survives, survived, withstands. "It survived every constraint you added" names no constraint, so the reader cannot check it. Say which constraint and what met it.

    Keep a word when it describes something other than the writer's sincerity. "A candid photo", "a frank disagreement between the two teams", "a straight line", and "straight quotes" (rule 20) are all fine. Outcome words are the exception: "the migration withstood the load test" is fine because the load test is named.

    When a claim needs credibility, supply the evidence instead of the label. "Two honest options:" becomes "Two options:". Add to this list whenever a word gets caught doing that job. Rule 47 handles the same move at sentence level, where a category noun or a spatial verb replaces the evidence.

### Occasion-built metaphor

42. **Occasion-built metaphor.** A figure of speech coined inside the sentence to relate two abstract things, where the image carries the claim instead of a stated relation. "The ladder, the roles and the triggers are machinery for moving that gate" tells the reader that three things serve one purpose and never says how. "The gate" is fine, because it names a real checkpoint with a real rule. "Machinery for moving it" is the failure, because nothing in the sentence is a machine.

    The tell is an abstract noun doing physical work in a sentence about a relation: machinery, engine, plumbing, spine, rails, seam, skeleton, backbone, glue, flywheel, the front desk, the bus. Established terms the reader already uses are fine, like "the bus" for a message broker or "a seam" for a module boundary. The problem is coining one on the spot for one paragraph.

    Two tests:
    - Delete the image. If the sentence still makes the point, the image was decoration.
    - Keep the image. If you cannot say item by item what it maps to, the relation is not worked out yet, and the reader will not work it out either.

    The fix is to state the relation. "The ladder changes exactly one variable: what tells the loop to stop" says what the machinery sentence was gesturing at.

    Judgement call: a metaphor works when it names something the reader owns and the rest of the sentence says something checkable about it. "That leaves the agent as the only writer on the branch" is concrete. "That leaves the agent as the flywheel" is not. When the image is load-bearing, the sentence is arguing with a picture instead of with a fact.

### Compression

Rules 28, 34, 36, 40 and 42 all expand a compressed phrase into an explanation. Rule 43 pulls the other way, and it is what keeps a rewrite inside the length budget in the process section.

43. **Over-compression.** Dropped articles, verbless fragments, arrows, symbols, and private abbreviations make the reader decode instead of read. "Parser rejects bad date → exit 2, no write" becomes "The parser rejects a bad date, exits with code 2, and writes nothing." "Retry cap 3, then surface" becomes "The client retries three times, then reports the error." Write whole sentences with their articles and verbs, and spell out arrows and abbreviations. This rule disagrees with rules 34 through 36 about one narrow case: clipped phrasing is a tell in prose, and telegraphic phrasing is a tell in commit messages and error strings.

44. **Error messages and code comments.** Name the thing that failed, the condition that caused it, and what the caller should do. "Error: invalid input" tells the reader nothing. "Field `due_date` must be an ISO date, got `tomorrow`" names the field, the rule, and the bad value. For comments, explain the reason the code exists, including the constraint or incident that forced it, and never restate what the line already says. A comment that repeats the code adds a second place to keep in sync.

### Mannered prose

45. **Mannered prose.** A metaphor or flourish where a literal phrase exists: personified code ("the plan holds it"), figurative verbs ("rides along", "stands on", "wants to"), and stock framing phrases. "A dial worth turning" becomes "a parameter worth varying". Aphorisms and rhetorical fragments belong to rules 34 through 36, and invented metaphors to rule 42. Say what you mean.

### Slogan labels and phantom locations

46. **Slogan labels.** A bold or headed equation of two nouns the reader already knows, followed by a body that describes something else. "**B. The channel is the session.** The channel persists, the conversation accumulates in it, and a mention starts a turn." The label claims sameness. The body lists duration, content, and a trigger, so the reader has to work out where the sameness went. The label only parses after the body, and the body never confirms it. Test: read the body alone and ask whether it establishes the label. If it does not, cut the label and make the relation the first sentence: "A mention opens a turn in the channel's conversation. The channel keeps that conversation, so a later mention can read what earlier turns said." A heading that names the claim its paragraph proves is fine, so keep the ones that pass the test.

    This is rule 17's bold lead-in gone wrong, and it usually arrives with rule 35's clipped predicates underneath.

47. **Diagnostic nouns and phantom locations.** Two constructions that assert a relation and skip it.

    The first puts a category in the subject where the instance belongs: "the tell is", "the signs are", "the criteria", "the signal", "the marker", "X tells you". "The tell is the triple" makes the reader work out which triple and why a triple indicates anything. Name the instance and its consequence: "three clauses in a row with no connector make the rhythm carry a claim the words never made".

    The second puts a non-physical thing at an address: "the risk sits in the retry loop", "the bug lives in the parser", "the answer sits in the config". A spatial verb promises a location, and the noun it governs has none. Say what happens: "the retry loop has no cap, so one 502 repeats until the account hits its rate limit".

    "Two smaller tells sit on top" is both patterns in one sentence, plus a count announced before the items. Fix the two nouns first, then delete the count. Rules 40 and 42 cover the neighbouring cases of a verb carrying the mechanism and a metaphor carrying the relation.

    This file talks about writing, so "tell" is a technical term here. It is fine when the thing named is an indicator the reader already uses, as in "the `->` arrow is a tell in commit messages". It is not fine when it names a category and leaves the reader to find the instance.
