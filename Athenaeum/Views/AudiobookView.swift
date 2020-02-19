/**
 AudiobookView.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import SwiftUI

struct AudiobookView: View {
    let book: Audiobook

    var body: some View {
        VStack {
            if book.getCover() != nil {
                Cover(data: book.getCover()!)
            }
            VStack(alignment: .leading) {
                if book.title.contains(":") {
                    ForEach(book.title.components(separatedBy: ":"), id: \.self) { title in
                        Text(title.trimmed)
                            .font(.headline)
                    }
                } else {
                    Text(book.title)
                        .font(.headline)
                }
                HStack {
                    Text(book.author)
                        .font(.subheadline)
                    if book.publicationDate != nil {
                        Spacer()
                        Text(book.publicationDate!)
                            .font(.subheadline)
                    }
                }
            }
            if book.summary != nil {
                Divider()
                ScrollView {
                    SummaryView(summary: book.summary!)
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
            .shadow(radius: 10)
    }
}

#if DEBUG
    struct AudiobookView_Previews: PreviewProvider {
        static var previews: some View {
            Group {
                AudiobookView(book: Audiobook(fromFileWithPath: "/Users/ckerson/Music/TWoK.m4b"))
            }
        }
    }
#endif
