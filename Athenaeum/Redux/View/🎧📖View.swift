/**
 🎧📖View.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Combine
import SwiftUI

struct AudioBookView: View {
    @ObservedObject var viewModel: AudiobookViewModel

    var cancelable: AnyCancellable?

    init(_ viewModel: AudiobookViewModel) {
        self.viewModel = viewModel

        self.cancelable = self.viewModel.objectWillChange.sink { _ in
            log.info("Updating \(AudioBookView.self)")
        }
    }

    var body: some View {
        VStack {
            Button(action: {
                self.viewModel.importFromOpenDialog()
            }) {
                Text("Import")
            }
            if viewModel.isImporting {
                Text("Importing...").bold()
            }
            ForEach(viewModel.audiobooks.sorted(by: {
                if let authorZero = $0.authors?.first, let authorOne = $1.authors?.first {
                    if authorZero != authorOne {
                        return authorZero.lastName < authorOne.lastName
                    }
                }
                if let pubDateZero = $0.publicationDate,
                    let pubDateOne = $1.publicationDate {
                    return pubDateZero < pubDateOne
                } else {
                    return $0.title < $1.title
                }
            }), id: \.id) { book in
                HStack {
                    Text(book.title)
                    Spacer()
                    Unwrap(book.getAuthorsString()) { author in
                        Text(author)
                    }
                }
            }
        }
        .padding()
        .frame(minWidth: 400)
    }
}

#if DEBUG
    let sampleDescription = """
    I long for the days before the Last Desolation.

    The age before the Heralds abandoned us and the Knights Radiant turned against us. A time when there was still magic in the world and honor in the hearts of men.

    The world became ours, and yet we lost it. Victory proved to be the greatest test of all. Or was that victory illusory? Did our enemies come to recognize that the harder they fought, the fiercer our resistance? Fire and hammer will forge steel into a weapon, but if you abandon your sword, it eventually rusts away.

    There are four whom we watch. The first is the surgeon, forced to forsake healing to fight in the most brutal war of our time. The second is the assassin, a murderer who weeps as he kills. The third is the liar, a young woman who wears a scholar's mantle over the heart of a thief. The last is the prince, a warlord whose eyes have opened to the ancient past as his thirst for battle wanes.

    The world can change. Surgebinding and Shardwielding can return; the magics of ancient days become ours again. These four people are key.

    One of them may redeem us. And one of them will destroy us.
    """

    let sampleAudiobook = AudioBook(id: UUID(),
                                    location: URL(fileURLWithPath: "/Users/ckerson/Music/TWoK.m4b"),
                                    contentsHash: "",
                                    title: "The Way of Kings",
                                    authors: [Author(firstName: "Brandon", lastName: "Sanderson")],
                                    narrator: "Michael Kramer & Kate Reading",
                                    publicationDate: "2010-08-31",
                                    isbn: "0765326353",
                                    bookDescription: sampleDescription,
                                    series: Series(title: "The Stormlight Archive",
                                                   entry: "1"))
    let sampleStore = Store<GlobalAppState>(reducer: appStateReducer,
                                            middleware: [logMiddleware],
                                            state: GlobalAppState(audiobookState: AudiobookState(audiobooks: [sampleAudiobook
                                                    .id: sampleAudiobook])))

    struct AudiobookListView_Previews: PreviewProvider {
        static var previews: some View {
            AudioBookView(AudiobookViewModel(store: sampleStore))
        }
    }

#endif
