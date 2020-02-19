/**
 AudiobookListView.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import SwiftUI

struct AudiobookListView: View {
    let books: [Audiobook]

    var body: some View {
        List(books) { book in
            AudiobookCellView(book: book)
        }
    }
}

struct AudiobookList_Previews: PreviewProvider {
    static var previews: some View {
        AudiobookListView(books: [
            Audiobook(fromFileWithPath: "/Users/ckerson/Music/TWoK.m4b"),
            Audiobook(fromFileWithPath: "/Users/ckerson/Music/In the Labyrinth of Drakes.m4b"),
            Audiobook(fromFileWithPath: "/Users/ckerson/Music/Smarter Faster Better.m4b"),
        ])
    }
}
