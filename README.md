# Allonsay!

[![Build Status](https://travis-ci.com/chunfeilung/allonsay.svg?branch=master)](https://travis-ci.com/chunfeilung/allonsay)

![Allonsay](allonsay.png)

In macOS you can select any piece of text anywhere, and have it read out
loud by the system’s voice.

This works if all content that you consume is in the same language, but falls
flat on its face when it has to deal with multiple languages: some words will
“merely” be incomprehensible, but others might even be skipped over completely.
Yes, you can manually select another voice that does fit the language, but this
can become a bit of a hassle if you often consume content in different
languages.

This is where Allonsay comes in: Allonsay is a tiny utility that enables macOS
to read text out loud using the correct language.

_In case you were wondering, Allonsay is pronounced as “Allons-y”, which means
“Let’s go” in French. Coincidentally, this is also why it was developed using
the Go programming language._

## Features
Allonsay is still under active development, so here’s a brief overview of the
things it can and cannot do.

### What works
* You can pass a text snippet to Allonsay in Dutch, English, and/or Chinese, and
  it’ll automatically figure out the best way to read it out to you… most of the
  time.
* Speech synthesis using macOS’ built-in `say` command, with the following
  voices:
  * Samantha (_English_)
  * Claire (_Dutch_)
  * Sin-ji (_Cantonese_)

### What doesn’t work (yet)
* Any language that isn’t Dutch, English, or Cantonese. Allonsay doesn’t know
  the difference between traditional and simplified Chinese yet, and if you give
  it a text that uses the western alphabet, it will attempt to read it using an
  English or Dutch voice.
* Even if you ask Allonsay to read a text snippet that _is_ in Dutch or English,
  it might still use the wrong voice; it currently classifies languages based on
  letter frequency, but this method doesn’t work very well on short text. In the
  future I may add a different classifier based on n-grams (or whatever works
  best).
* Ability for users to configure which languages and voices they want to use.

## Getting started
Before you actually get started, let me remind you that Allonsay is:
* Under active development and not even close to something that I’d be willing
  to call “1.0”. Having said that, it’s already in reasonably usable state;
* A hobby project that’ll probably be abandoned the moment it works well enough
  for me;
* Available under [an MIT license](LICENSE.md). Unless you have a technical or
  legal background, that basically means that you can use Allonsay for free
  – all you need is an expensive Mac.

### Prerequisites
First, make sure that the voices Samantha, Claire, and Sin-ji are installed on
your Mac. You can verify this by navigating to
_System Preferences > Accessibility > Speech > System Voice > Customise_ and
install them there, if needed.

### Installation
Allonsay is distributed as a command-line utility. Download the latest _binary_
release from the “Releases” page, make sure it’s executable, and put it in
`/usr/local/bin` (or any other folder to your liking).

### Usage (CLI)
Pass a string as the first (and only) parameter; if the string contains spaces,
it should be enclosed using straight quotation marks:

```
allonsay supercalifragilisticexpialidocious
allonsay hottentottententententoonstelling
allonsay 氏時時適市視獅
allonsay "HKIA is colloquially known as 赤鱲角機場 in Cantonese"
allonsay "In het Kantonees heet Nederland 荷蘭"
```

### Usage (macOS integration)

Your Mac should come preinstalled with _Automator_.

1. Open _Automator_ and choose “New Document”.
2. _Automator_ will prompt you to choose a type for your document. Select
   “Service” by double-clicking it.
3. You now see a window with actions on the left, and an empty canvas on the
   right.
4. Look for “Run Shell Script” in the list of actions on the left, and drag it
   to the canvas.
5. Make sure the input is passed as arguments (on the right). Enter the
   following in the big textarea:
   ```
   /usr/local/bin/allonsay "$*"
   ```
   Replace the path to `allonsay` if you chose a different location.

6. Save the service using a name like “Read this text”

Many macOS applications will now show a new “Read this text” entry in the
context menu whenever you have selected some text.

## Contact
If you have any questions, issues, or requests, feel free to:
 * [create an issue](https://github.com/chunfeilung/allonsay/issues/new)
 * drop [the author](https://chuniversiteit.nl) a message via
   [Twitter](https://twitter.com/chunfeilung),
   [Facebook](https://www.facebook.com/lungchunfei),
or [LinkedIn](https://linkedin.com/in/chunfeilung/).
