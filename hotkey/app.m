// Copyright 2020 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

#import <Foundation/Foundation.h>
#import <Cocoa/Cocoa.h>
#import <Carbon/Carbon.h>

extern void
go_hotkey_callback(void* handler);

static OSStatus
_hotkey_handler(EventHandlerCallRef nextHandler, EventRef theEvent, void *userData) {
	EventHotKeyID k;
	void *go_hotkey_handler = userData;

	GetEventParameter(theEvent, kEventParamDirectObject, typeEventHotKeyID, NULL, sizeof(k), NULL, &k);
	if (k.id == 1) go_hotkey_callback(go_hotkey_handler);
	return noErr;
}

int
register_hotkey(void* go_hotkey_handler) {
	EventHotKeyID hotKeyID  = {'htk1', 1};
	EventTypeSpec eventType = {kEventClassKeyboard, kEventHotKeyPressed};
	InstallApplicationEventHandler(&_hotkey_handler, 1, &eventType, go_hotkey_handler, NULL);
	EventHotKeyRef hotKeyRef;
	OSStatus s = RegisterEventHotKey(1, controlKey+optionKey, hotKeyID,
		GetApplicationEventTarget(), 0, &hotKeyRef);
	if (s != noErr) {
		return -1;
	}
	return 0;
}

void
app_run() {
	[NSApplication sharedApplication];
	[NSApp disableRelaunchOnLogin];
	[NSApp run];
}