/**
 LibraryView.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import SwiftUI

struct LibraryView<R>: View where R: Repository, R.EntityObject: Audiobook {
    @ObservedObject var repository: R

    init(withRepository repository: R) {
        self.repository = repository
    }

    var body: some View {
        NavigationView {
            List {
                Section(header: HeaderView()) {
                    ForEachAudiobookView(inRepository: repository)
                }
            }
            .frame(minWidth: 425, maxWidth: 425)
        }
        .listStyle(SidebarListStyle())
        .frame(minWidth: 850,
               maxWidth: 850,
               minHeight: 400,
               maxHeight: .infinity)
    }
}

struct HeaderView: View {
    var body: some View {
        VStack {
            HStack(spacing: 20) {
                Text("Library")
                    .layoutPriority(1)
                    .font(.largeTitle)
                Spacer()
            }
        }.padding(.bottom)
    }
}

#if DEBUG
    struct LibraryView_Previews: PreviewProvider {
        static var previews: some View {
            Group {
                LibraryView(withRepository: MockRepo())
                    .environment(\.colorScheme, .light)
                LibraryView(withRepository: MockRepo())
                    .environment(\.colorScheme, .dark)
            }
        }
    }
#endif
