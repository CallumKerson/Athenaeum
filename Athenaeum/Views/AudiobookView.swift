/**
 AudiobookView.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import SwiftUI

struct AudiobookView: View {
    let book: Audiobook

    var body: some View {
        VStack {
            Unwrap(book.getCover()) { cover in
                Cover(data: cover)
            }
            VStack(alignment: .center) {
                if book.title.contains(":") {
                    ForEach(book.title.components(separatedBy: ":"),
                            id: \.self) { title in
                        Text(title.trimmed)
                            .font(.headline)
                    }
                } else {
                    Text(book.title)
                        .font(.headline)
                }
                Unwrap(book.series) { series in
                    Text("Book \(series.entry) of \(series.title)")
                }
            }
            VStack(alignment: .leading) {
                HStack {
                    Text(book.author)
                        .font(.subheadline)
                    Unwrap(book.publicationDate) { date in
                        Spacer()
                        Text(date)
                            .font(.subheadline)
                    }
                }
            }
            Unwrap(book.narrator) { narrator in
                VStack(alignment: .leading) {
                    HStack {
                        Text("Narrated by \(narrator)")
                        Spacer()
                    }
                }
            }
            Unwrap(book.summary) { summary in
                Divider()
                ScrollView {
                    SummaryView(summary: summary)
                }
                .frame(minHeight: 50)
            }
        }
        .frame(idealWidth: 200, maxWidth: 400)
        .padding()
    }
}

struct Cover: View {
    let data: Data

    var body: some View {
        Image(nsImage: NSImage(data: data)!)
            .resizable()
            .scaledToFit()
            .frame(width: 400, height: 400)
            .clipShape(RoundedRectangle(cornerRadius: 5))
            .shadow(radius: 5)
    }
}

#if DEBUG
    struct AudiobookView_Previews: PreviewProvider {
        static var previews: some View {
            Group {
                AudiobookView(book: previewAudiobooks[0])
                AudiobookView(book: previewAudiobooks[1])
                    .environment(\.colorScheme, .dark)
                AudiobookView(book: previewAudiobooks[3])
            }
        }
    }
#endif
