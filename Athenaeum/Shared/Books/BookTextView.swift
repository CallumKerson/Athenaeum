/**
 BookTextView.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import SwiftUI

struct BookTextView: View {
    let metadata: BookMetadata

    var body: some View {
        VStack(alignment: HorizontalAlignment.leading) {
            VStack(alignment: HorizontalAlignment.leading) {
                Text(metadata.title).font(.title)
                if let series = metadata.series {
                    Text("Book \(series.entry.asString) of \(series.title)")
                }
            }.padding(.bottom)

            if let author = (metadata.authors?.author) {
                Text("Written by ") + Text(author)
            }
            if let releaseDate = (metadata.publicationDate?.getDateAsString) {
                Text("Released on ") + Text(releaseDate)
            }
            if let narrator = (metadata.narrators?.author) {
                Text("Narrated by ") + Text(narrator)
            }
            if let isbn = metadata.isbn {
                Text("ISBN ") + Text(isbn)
            }
        }
        .font(.subheadline)
    }
}

struct BookTextView_Previews: PreviewProvider {
    static var previews: some View {
        var sampleMetadata = BookMetadata(title: "Murder on the Orient Express")
        sampleMetadata.authors = ["Agatha Christie"]
        sampleMetadata.narrators = ["David Suchet"]
        return BookTextView(metadata: sampleMetadata)
    }
}
