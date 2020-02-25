/**
 AudiobookListView.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import SwiftUI

struct ForEachAudiobookView<Lib>: View where Lib: Library {
    @ObservedObject var library: Lib

    init(inLibrary library: Lib) {
        self.library = library
    }

    var body: some View {
        ForEach(library.🎧📚.sorted(by: {
            if $0.author != $1.author {
                return $0.author < $1.author
            } else if let pubDateZero = $0.publicationDate,
                let pubDateOne = $1.publicationDate {
                return pubDateZero < pubDateOne
            } else {
                return $0.title < $1.title
            }
            }), id: \.title) { book in
            NavigationLink(destination: AudiobookView(book: book)) {
                AudiobookCellView(book: book)
            }
        }
    }
}

struct AudiobookList_Previews: PreviewProvider {
    static var previews: some View {
        ForEachAudiobookView(inLibrary: MockLibrary())
    }
}
