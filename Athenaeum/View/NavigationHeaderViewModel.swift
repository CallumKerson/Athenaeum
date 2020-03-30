/**
 NavigationHeaderViewModel.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Cocoa
import Combine
import Foundation
import SwiftUI

class NavigationHeaderViewModel: ObservableObject {
    var isImporting: Bool = false
    var isErrored: Bool = false
    let objectWillChange = ObservableObjectPublisher()

    private var didStateChangeCancellable: AnyCancellable?
    var store: Store<GlobalAppState>

    init(store: Store<GlobalAppState>) {
        self.store = store
        self.didStateChangeCancellable = self.store.stateSubject.sink(receiveValue: {
            let incomingImportState = ($0.audiobookState.audiobooks.values.filter { $0.isLoading }
                .count != 0)
            if self.isImporting != incomingImportState {
                self.isImporting = incomingImportState
                self.objectWillChange.send()
            }
            let incomingErroredState = ($0.audiobookState.audiobooks.values.filter { $0.isErrored }
                .count > 0)
            if self.isErrored != incomingErroredState {
                self.isErrored = incomingErroredState
                self.objectWillChange.send()
            }
        })
    }

    func importButtonAction() {
        importFromOpenDialog(store: self.store)
    }

//    func errorIconHoverAction() {
//        let popover = NSPopover()
//        popover.contentSize = NSSize(width: 400, height: 500)
//        popover.behavior = .transient
//        popover.contentViewController = NSHostingController(rootView: ErrorsPopoverView())
//        popover.show(relativeTo: button.bounds, of: button, preferredEdge: NSRectEdge.minY)
//    }
}
