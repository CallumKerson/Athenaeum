/**
 SummaryView.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import SwiftUI

struct SummaryView: View {
    let lines: [String]

    init(summary: String) {
        lines = summary.components(separatedBy: "<br />")
    }

    var body: some View {
        VStack(alignment: .leading) {
            ForEach(lines, id: \.self) { line in
                self.getLine(line: line)
            }
        }
        .fixedSize(horizontal: false, vertical: true)
    }

    func getLine(line: String) -> Text {
        if line.hasPrefix("<i>") && line.hasSuffix("</i>") {
            return Text(line.replacingOccurrences(of: "<i>", with: "")
                .replacingOccurrences(of: "</i>", with: ""))
                .italic()
        }
        if line.hasPrefix("<b>") && line.hasSuffix("</b>") {
            return Text(line.replacingOccurrences(of: "<b>", with: "")
                .replacingOccurrences(of: "</b>", with: ""))
                .bold()
        }
        return Text(line
            .replacingOccurrences(of: "<b>", with: "")
            .replacingOccurrences(of: "</b>", with: "")
            .replacingOccurrences(of: "<i>", with: "")
            .replacingOccurrences(of: "</i>", with: ""))
    }
}

struct SummaryView_Previews: PreviewProvider {
    static var previews: some View {
        SummaryView(summary: "<i>I long for the days before the Last Desolation.<br /><br />The age before the Heralds abandoned us and the Knights Radiant turned against us. A time when there was still magic in the world and honor in the hearts of men.<br /><br />The world became ours, and yet we lost it. Victory proved to be the greatest test of all. Or was that victory illusory? Did our enemies come to recognize that the harder they fought, the fiercer our resistance? Fire and hammer will forge steel into a weapon, but if you abandon your sword, it eventually rusts away.<br /><br />There are four whom we watch. The first is the surgeon, forced to forsake healing to fight in the most brutal war of our time. The second is the assassin, a murderer who weeps as he kills. The third is the liar, a young woman who wears a scholar's mantle over the heart of a thief. The last is the prince, a warlord whose eyes have opened to the ancient past as his thirst for battle wanes.<br /><br />The world can change. Surgebinding and Shardwielding can return; the magics of ancient days become ours again. These four people are key.<br /><br />One of them may redeem us. And one of them will destroy us.</i><br /><br />From Brandon Sanderson-who completed Robert Jordan's The Wheel of Time-comes The Stormlight Archive, an ambitious new fantasy epic in a unique, richly imagined setting. Roshar is a world relentlessly blasted by awesome tempests, where emotions take on physical form, and terrible secrets hide deep beneath the rocky landscape.<br /><br /><i>Speak again the ancient oaths</i><br /><b>Life before death. Strength before weakness. Journey before destination.</b><br /><i>and return to men the Shards they once bore. The Knights Radiant must stand again!</i>")
    }
}
