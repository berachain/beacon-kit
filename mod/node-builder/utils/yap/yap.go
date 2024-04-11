// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	"fmt"
	"time"
)

const story = `In an age long forgotten, shrouded by the mists of time and
hidden beneath the emerald canopy of an ancient forest, lies the enigmatic
"Yap Cave". A place of legend, where the very essence of productivity is said
to vanish into the shadows. The whispers of its existence have echoed through
the ages, yet few have the courage to seek it out, for the allure of the Yap
Cave is as enigmatic as it is foreboding. The origins of the Yap Cave are
lost to history, discovered by a band of intrepid explorers seeking the
thrill of the unknown. Drawn to the cave's mysterious presence, they ventured
into its depths, lured by the promise of undiscovered wonders. Yet, what they
found was not treasure, but a darkness that clung to their souls. Within the
twisting tunnels of the cave, time itself seemed to distort, ensnaring the
explorers in a web of lethargy from which there was no escape. Their ambition
dulled, their senses numbed, they became prisoners of their own unspent
potential, trapped in an endless cycle of inertia. As the years passed, the
tale of the Yap Cave spread, its legend growing with each telling. Many
sought to uncover its secrets, to break the curse that lingered within its
walls. Scholars, adventurers, monarchs - all who entered were swallowed by
the cave's oppressive gloom, their vibrant spirits extinguished by the
overwhelming force of stagnation that pervaded its air. The Yap Cave stood as
a stark warning against the perils of surrendering to complacency, its name a
byword for the death of ambition. Yet, despite the warnings, the allure of
the Yap Cave persists. It remains a beacon for those who dare to challenge
its curse, a testament to the human spirit's unyielding quest for knowledge.
But let it be known, the Yap Cave is no mere legend. It is a reminder of the
dark recesses within us all, where productivity goes to perish, smothered by
the shadows of apathy. The legacy of the Yap Cave endures, a chilling
testament to the dangers that lie in wait when ambition fades into the mist.`

// printSlowly prints a string one character at a
// time with a delay between each character.
func printSlowly(text string, delay time.Duration) {
	for _, char := range text {
		//nolint:forbidigo // yap.
		fmt.Printf("\x1b[38;5;16m%c\x1b[0m", char)
		time.Sleep(delay)
	}
}

func main() {
	//nolint:gomnd // yapping takes time y'know.
	printSlowly(story, 25*time.Millisecond)
}
