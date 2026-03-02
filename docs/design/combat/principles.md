# The 10 Design Principles

**Author**: Karsten Huttelmaier — co-authored with Claude

These aren't guidelines. They're laws. When we're stuck on an implementation decision, we come back here.

---

## Principle 1: "30-90 Second Fights"

*The Covert Action Firewall.*

**What it means:** No single combat encounter should take more than 90 seconds for a competent player. Boss fights can stretch to 2 minutes max. If a fight goes longer, something is wrong.

**Why it matters:** Meier learned from Covert Action that two competing games destroy each other. Our dungeon exploration and platformer combat are different games. The ONLY reason they work together is that combat is SHORT enough to feel like a punctuation mark in the dungeon journey — an exclamation point, not a paragraph.

**The scenario that proves it:** Player has been exploring for 3 minutes. Found the key. Knows where the locked door is. Walks into a room with an enemy. Fights for 45 seconds. Wins. Gets a health drop. Returns to the dungeon and immediately heads for the door. The combat was an EVENT in their exploration story, not a detour from it.

**The scenario that violates it:** Same player, same situation. The fight drags on for 4 minutes. By the time it's over, they've forgotten which direction the locked door was. They have to re-orient. The dungeon exploration feels like it restarted. The two games fought each other.

**The test:** Time every fight during playtesting. If average fight time exceeds 60 seconds, cut enemy HP in half. If it exceeds 90 seconds, cut it in half AGAIN. Meier: "Double it or cut it in half."

---

## Principle 2: "HP Is the Universal Currency"

*The thread that stitches two games into one.*

**What it means:** The player has one HP pool. It carries from room to room, from dungeon to combat and back. There is no "combat HP" and "exploration HP." There is just HP. It is the single shared resource that makes the dungeon and combat into one unified game.

**Why it matters:** Without shared HP, the dungeon and combat are disconnected experiences. With it, every combat encounter has STAKES beyond the fight itself. Taking 2 damage in Fight #2 affects whether you can survive Fight #4. This is what transforms a series of isolated fights into a CAMPAIGN.

**The scenario that proves it:** Player finishes a fight with 2 HP remaining. The next room might have a health potion — or another enemy. They have to decide: explore cautiously or beeline for the exit? This tension doesn't exist if HP resets between fights.

**The scenario that violates it:** HP resets to full after each fight. Now every fight is identical — no stakes, no resource management, no connection between encounters. Each fight is a standalone minigame. The dungeon becomes a hallway between disconnected combat puzzles. We've made Covert Action.

**The test:** If a player says "that fight didn't matter" or "I don't care about taking damage," something is broken. Damage must cost something BEYOND the fight.

---

## Principle 3: "Every Room Is a Decision"

*No filler rooms.*

**What it means:** Every room the player enters should present at least one meaningful choice: pick up an item (but it takes an inventory slot), fight an enemy (but costs HP), take a shortcut (but miss potential loot), use a health potion (but it's gone forever).

**Why it matters:** Meier: "A game is a series of interesting decisions." A room with nothing to do, nothing to decide, and nothing at stake is dead space. The player walks through it and feels nothing. That's a missed opportunity.

**The scenario that proves it:** Player enters a room. Sees a health potion AND a sword. They already have 3 items and a full inventory. They have to DROP something to pick either one up. Which is more valuable right now? That depends on their HP, upcoming enemies, and playstyle. THAT'S a room.

**The scenario that violates it:** Player enters 3 empty rooms in a row. No items, no enemies, no choices. They're just pressing 'd' to walk east. They're not playing the game — they're commuting through it.

**The test:** Walk through the dungeon and note every room where the player presses a key without thinking. Those rooms need something — an item, an enemy, a fork in the path, a clue.

**Generator implication:** The generator must ensure no more than 1 empty room in a row. Every 2nd or 3rd room should have SOMETHING (item, enemy, or branching path).

---

## Principle 4: "The Player Is the Star"

*Meier: "Our job is to impress you with yourself."*

**What it means:** The game should make the player feel skillful, powerful, and clever. Not by being easy — by being GENEROUS in how it interprets the player's actions, and CLEAR in how it communicates challenges.

**In combat, this means:**
- **Coyote time** (4 frames): Player walks off a ledge? They can still jump for 4 frames. Instead of "I fell because I was 1 pixel late," it becomes "Wow, I saved that jump at the last second!" The player takes credit. The game does the work.
- **Jump buffering** (4 frames): Player presses jump 4 frames before landing? The jump happens on landing. Instead of "the game ate my input," it becomes "I timed that perfectly."
- **Attack hitboxes bigger than visual**: The nail slash hits slightly further than it looks. Near-misses become hits. The player feels precise.
- **Enemy telegraphs**: Every attack has a 0.3-0.5s windup. The player can ALWAYS see it coming. When they dodge, they feel smart. When they get hit, they think "I saw that coming and reacted too slow" — not "that was unfair."

**In the dungeon, this means:**
- Clear room descriptions that give useful information
- The map always shows where you've been
- Items have obvious names that suggest their use
- The key is never behind the locked door (already enforced by the generator!)

**The test:** After every death, ask: "Did the player feel like it was THEIR fault?" If they blame the game, something is unfair. If they blame themselves, the design is working. Meier: "Players take credit for wins, blame the game for losses."

---

## Principle 5: "The Contrast IS the Spectacle"

*Don't upgrade the dungeon. Keep it quiet. The louder combat is, the louder the contrast.*

**What it means:** The dungeon exploration should stay small, monochrome, understated — the box-drawing map, text descriptions, simple commands. Do NOT add color, animation, or visual flair to the dungeon. The magic happens when a tiny, quiet text game EXPLODES into full-screen pixel art combat. The bigger the gap, the bigger the wow.

**Think of it like a movie:** The quiet scenes make the action scenes hit harder. If every scene is an explosion, nothing is exciting. If the hero is sitting quietly in a room and then suddenly a wall explodes — THAT'S cinema.

**The scenario that proves it:** Someone watches over a player's shoulder. They see a little box-drawing map. "Oh, a text adventure." Then the player walks into a room and the screen TRANSFORMS into full-color pixel art with a warrior fighting a monster. "Wait, WHAT?!" That reaction is the feature.

**The scenario that violates it:** We add color, animation, and visual polish to the dungeon map. Now the dungeon is pretty. But when combat triggers, the mode switch feels like a minor visual upgrade instead of a jaw-dropping transformation. The "WHAT?!" becomes "oh, cool."

**The test:** Show the game to someone who hasn't seen it before. Watch their face when combat triggers for the first time. If they don't react with surprise, the contrast isn't big enough.

---

## Principle 6: "Simple Systems, Complex Interactions"

*Meier: "Design small and simple systems that will interact in a complex manner."*

**What it means:** Each system in the game should be explainable in one sentence:
- "HP goes down when you get hit and persists between rooms."
- "You can carry up to 4 items."
- "Enemies drop loot when defeated."
- "The key opens the locked door."

None of these are complex. But when they INTERACT:
- "I have 2 HP, a key, a sword, and a health potion. The next room has an enemy guarding a shortcut to the locked door. If I fight and win, I save 3 rooms of walking. If I lose, game over. I could drink the potion for safety, but what if there's another fight? I could drop the sword to pick up the shield I saw earlier, but then the fight takes longer and I might take more hits..."

THAT'S complexity. And it emerged entirely from simple systems interacting.

**The test:** Can you explain every system to a 10-year-old in one sentence? If not, simplify. Can a sophisticated player find themselves agonizing over a decision for 10 seconds? If not, the interactions aren't rich enough.

**What NOT to do:** Don't add elemental damage types, stat buffs, weapon durability, hunger systems, or ability cooldowns. Each sounds like "just one more system" but each adds complexity to EVERY interaction. A game with 4 simple systems has 6 possible interactions (4 choose 2). A game with 8 systems has 28. Keep it tight.

---

## Principle 7: "Double It or Cut It in Half"

*Meier's tuning philosophy. The most practical rule on this list.*

**What it means:** When something doesn't feel right, don't nudge. LURCH. If enemies feel too tough, cut their HP in half. If attacks feel weak, double the damage. If the screen shake is wimpy, double it. If the fight is too long, halve the arena size.

**Why:** Small changes (5-10%) waste playtesting iterations because you can't FEEL the difference. You play, it seems about the same, you adjust again, play again... 7 iterations later you've moved 35% and wasted a day. One big change immediately tells you: "Is this even the right variable to adjust?"

**Application to our game:**
- First playtest: fights feel too long. **Cut enemy HP in half.** Don't reduce by 15%.
- Second playtest: fights are now too easy. **Double enemy attack damage.** Now fights are short AND scary.
- Third playtest: screen shake is barely noticeable. **Double the shake magnitude.** Ahh, NOW it feels like something.

**The test:** If you made a change and can't immediately feel the difference in the next playthrough, the change was too small.

---

## Principle 8: "Prototype the Fun First"

*Get to playable in the ugliest possible way. Then polish.*

**What it means:** The very first version of combat should be:
- A colored rectangle (the player) on a black screen
- Another colored rectangle (the enemy)
- Gravity. Jumping. Moving left/right.
- Pressing a key makes a rectangle appear next to the player (the attack hitbox)
- When the hitbox overlaps the enemy, the enemy rectangle blinks

No sprites. No particles. No screen shake. No transitions. No polish. Just rectangles.

If moving a rectangle around and hitting another rectangle feels FUN — if the weight of the jump feels right, if the attack timing feels snappy, if the collision feels fair — then we know the foundation is solid and everything else is polish.

If rectangles hitting rectangles is NOT fun, no amount of beautiful pixel art will save it.

**Meier: "There's just the hard, consistent work of making something a little better each day."** Start ugly. Make it play well. Then make it pretty.

**The milestone:** "I can't stop playtesting the rectangle version." When that happens, we know we've found the fun.

---

## Principle 9: "One Good Game"

*The Covert Action Rule, restated as a positive.*

**What it means:** This is a dungeon exploration game with a combat mechanic. It is NOT a dungeon game AND a platformer game. The platformer combat SERVES the dungeon. The dungeon is the meal; combat is the seasoning.

**How to test this:** Ask: "If I removed combat entirely, would there still be a game?" The answer must be YES. The dungeon with its keys, locks, items, and exploration must stand on its own (it already does!). Combat enhances it — adds stakes, variety, spectacle — but doesn't replace it.

**How to test the inverse:** "If I removed the dungeon entirely, would the combat be a good standalone game?" The answer should be NO. The combat is too short, too simple, and too repetitive to stand alone. It needs the dungeon to give it context, stakes, and variety. THIS IS CORRECT. Combat that could stand alone would be competing with the dungeon (Covert Action).

**The scenario that violates it:** We get excited about combat and keep adding features — complex combos, skill trees, equipment management, multiple enemy phases. Suddenly combat is a 5-minute ordeal with its own meta-game. Players start playing FOR the combat and viewing the dungeon as an annoying hallway between fights. We've accidentally made two games.

---

## Principle 10: "One More Room"

*Everything exists to make the player open that next door.*

**What it means:** This is the ultimate test. After every action the player takes — exploring a room, fighting an enemy, picking up an item — they should feel pulled toward the next room. "What's in there? Can I handle it? Is it the key? Is it a health potion? Is it something I've never seen?"

**The three motivators that keep them moving:**
1. **Curiosity**: "What's behind that door?" (the unknown)
2. **Greed**: "That enemy might drop something good." (the reward)
3. **Momentum**: "I'm already halfway through, might as well keep going." (sunk cost, but in a FUN way)

**The moment we've succeeded:** A player finishes a fight with 1 HP. Logically, they should retreat or play safe. Instead, they say "but the next room might have a health potion..." and they push forward. THAT'S "one more room."

**The moment we've failed:** A player finishes a fight and thinks "ugh, I hope I don't have to do that again" and quits. The combat was a chore, not a pull.

**The test:** Play the game yourself. Note the moment you want to stop. WHY did you want to stop? Fix that thing. Then note the moment you think "one more room." What caused that feeling? Amplify it.
