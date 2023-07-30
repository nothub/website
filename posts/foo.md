---
title: "Format Example"
description: "Foo text bla..."
date: "2022-10-29"
tags: ["foo"]
draft: true
---

These are the voyages of the Starship Enterprise. Its continuing mission, to explore strange new worlds, to seek out new
life and new civilizations, to boldly go where no one has gone before. We need to neutralize the homing signal. Each
unit has total environmental control, gravity, temperature, atmosphere, light, in a protective field. Sensors show
energy readings in your area. We had a forced chamber explosion in the resonator coil. Field strength has increased by
3,000 percent.

## Some text

Shields up. I recommend we transfer power to phasers and arm the photon torpedoes. Something strange on the detector
circuit. The weapons must have disrupted our communicators. You saw something as tasty as meat, but inorganically
materialized out of patterns used by our transporters. Captain, the most elementary and valuable statement in science,
the beginning of wisdom, is `I do not know`. All transporters off.

## More text

Sensors indicate no shuttle or other ships in this sector. According to coordinates, we have travelled 7,000 light years
and are located near the system J-25. Tractor beam released, sir. Force field maintaining our hull integrity. Damage
report? Sections 27, 28 and 29 on decks four, five and six destroyed. Without our shields, at this range it is probable
a photon detonation could destroy the Enterprise.

# Bar

## A code block with text

```
  /7/Y/^\\
  vuVV|C)|                           __,_
    \\|^ /                         .'  Y '>,
    )| \\)                        / _   _   \\
   //)|\\\\                       )(_) (_)(|}
  / ^| \\ \\                      {  4A   } /
 //^| || \\\\                      \\uLuJJ/\\l
| \"\"\"\"\"  7/>l__ _____ __       /nnm_n//
L>_   _-&lt; 7/|_-__,__-)\\,__)(\".  \\_>-&lt;_/D
)D\" Y \"c)  9)//V       \\_\"-._.__G G_c_.-jjs&lt;\"/ ( \\
 | | |  |(|               &lt; \"-._\"> _.G_.___)\\   \\7\\
  \\\"=\" // |              (,\"-.__.|\\ \\&lt;.__.-\" )   \\ \\
   '---'  |              |,\"-.__\"| \\!\"-.__.-\".\\   \\ \\
     |_;._/              (_\"-.__\"'\\ \\\"-.__.-\".|    \\_\\
     )(\" V                \\\"-.__\"'|\\ \\-.__.-\".)     \\ \\
        (                  \"-.__'\"\\_\\ \\.__.-\"./      \\ l
         )                  \".__\"\">>G\\ \\__.-\">        V )

It indicates a synchronic distortion in the areas emanating triolic waves. The cerebellum, the cerebral cortex, the brain stem, the entire nervous system has been depleted of electrochemical energy. Any device like that would produce high levels of triolic waves. These walls have undergone some kind of selective molecular polarization. I haven't determined if our phaser energy can generate a stable field. We could alter the photons with phase discriminators.
```

text above is in a `code` block without any highlighting

## A code block with code

text in a `code` block with `bash` highlighting

```bash
#!/usr/bin/env bash

set -eo pipefail

rand() {
    if [[ -n $2 ]]; then count=$2; else count=1; fi
    echo -n "$1" | grep -Eo '\S{1}' | shuf | head --lines "$count"
}

buf+="$(rand "☀☄🌎🌑🚀🛰🛸" 3)"
buf+="$(for _ in {1..7}; do rand ",;'~*°✦⊚⊙⨀⋇"; done)"
buf+="$(for _ in {1..20}; do rand ".⋅∙⋆"; done)"
buf+="$(for _ in {1..750}; do echo -n " "; done)"

echo "${buf}" | grep -Eo '[^\n]{1}' | shuf | tr -d '\n' | grep -Eo '.{60}'
```

### Inline Image

Unidentified vessel travelling at sub warp speed, bearing 235.7. Fluctuations in energy readings from it, Captain. All
transporters off. A strange set-up, but I'd say the graviton generator is depolarized. The dark colourings of the
scrapes are the leavings of natural rubber, a type of non-conductive sole used by researchers experimenting with
electricity. The molecules must have been partly de-phased by the anyon beam.

![unix spellbook](/assets/man.gif)
Ye olde unix spellbook

Exceeding reaction chamber thermal limit. We have begun power-supply calibration. Force fields have been established on
all turbo lifts and crawlways. Computer, run a level-two diagnostic on warp-drive systems. Antimatter containment
positive. Warp drive within normal parameters. I read an ion trail characteristic of a freighter escape pod. The bomb
had a molecular-decay detonator. Detecting some unusual fluctuations in subspace frequencies. Communication is not
possible. The shuttle has no power. ![inline image](/assets/chaos.webp) Using the gravitational pull of a star to slingshot
back in time? We are going to Starbase Montgomery for Engineering consultations prompted by minor read-out anomalies.
Probes have recorded unusual levels of geological activity in all five planetary systems. Assemble a team. Look at
records of the Drema quadrant. Would these scans detect artificial transmissions as well as natural signals?

### Figure image

Now what are the possibilities of warp drive? Cmdr Riker's nervous system has been invaded by an unknown microorganism.
The organisms fuse to the nerve, intertwining at the molecular level. That's why the transporter's biofilters couldn't
extract it. The vertex waves show a K-complex corresponding to an REM state. The engineering section's critical.
Destruction is imminent. Their robes contain ultritium, highly explosive, virtually undetectable by your transporter.

![sharks](/assets/sharks.gif)
Some text about hacker sharks.

It indicates a synchronic distortion in the areas emanating triolic waves. The cerebellum, the cerebral cortex, the
brain stem, the entire nervous system has been depleted of electrochemical energy. Any device like that would produce
high levels of triolic waves. These walls have undergone some kind of selective molecular polarization. I haven't
determined if our phaser energy can generate a stable field. We could alter the photons with phase discriminators.

# Dumm

## Lists

### Ordered List

1. asdf
2. bsdf
3. csdf

### Unordered List

- asdf
- bsdf
- csdf

## Quote

Now what are the possibilities of warp drive?

> Cmdr Riker's nervous system has been invaded by an unknown microorganism. The organisms fuse to the nerve,
> intertwining at the molecular level. That's why the transporter's biofilters couldn't extract it. The vertex waves
> show
> a K-complex corresponding to an REM state. The engineering section's critical. Destruction is imminent. ~foo bar

Their robes contain ultritium, highly explosive, virtually undetectable by your transporter.

## Tables

| Item         | Price    | # In stock |
|--------------|----------|------------|
| Juicy Apples | 1.99     | *7*        |
| Bananas      | **1.89** | 5234       |
| Shits A      | **2.39** | 9001       |
| Shits B      | **3.49** | 9002       |
| Shits C      | **4.59** | 9003       |
