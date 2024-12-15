// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

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
	//nolint:mnd // yapping takes time y'know.
	printSlowly(story, 25*time.Millisecond)
}
