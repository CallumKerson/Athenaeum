/**
 AudiobookListView.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import SwiftUI

struct ForEachAudiobookView<R>: View where R: Repository, R.EntityObject: Audiobook {
    @ObservedObject var repository: R

    init(inRepository repository: R) {
        self.repository = repository
    }

    var body: some View {
        ForEach(repository.items.sorted(by: {
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

#if DEBUG
    struct AudiobookList_Previews: PreviewProvider {
        static var previews: some View {
            ForEachAudiobookView(inRepository: MockRepo())
        }
    }
#endif
