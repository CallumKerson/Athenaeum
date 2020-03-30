/**
 SummaryView.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import SwiftUI

struct SummaryView: View {
    let lines: [String]

    init(summary: String) {
        self.lines = summary.components(separatedBy: "<br />")
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

#if DEBUG
    struct SummaryView_Previews: PreviewProvider {
        static var previews: some View {
            SummaryView(summary: sampleAudiobook.metadata!.summary!)
        }
    }
#endif
