/**
 AudiobookListView.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import SwiftUI

struct AudiobookListView: View {
    let books: [Audiobook]

    var body: some View {
        List(books, id: \.title) { book in
            AudiobookCellView(book: book)
        }
    }
}

struct AudiobookList_Previews: PreviewProvider {
    static var previews: some View {
        AudiobookListView(books: previewAudiobooks)
    }
}
