// Copyright 2018 visualfc. All rights reserved.

package tk

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Event struct {
	//The type field from the event. Valid for all event types.
	Type int

	//The send_event field from the event. Valid for all event types.
	//0 indicates that this is a “normal” event, 1 indicates that it is a “synthetic” event generated by SendEvent.
	Synthetic bool

	//The path name of the window to which the event was reported (the window field from the event). Valid for all event types.
	Widget Widget

	//The time field from the event. This is the X server timestamp (typically the time since the last server reset) in milliseconds,
	//when the event occurred. Valid for most events.
	Timestamp int64

	//The number of the button that was pressed or released.
	//Valid only for ButtonPress and ButtonRelease events.
	MouseButton int

	//The x and y fields from the event. For ButtonPress, ButtonRelease, Motion, KeyPress, KeyRelease, and MouseWheel events,
	//%x and %y indicate the position of the mouse pointer relative to the receiving window.
	//For Enter and Leave events, the position where the mouse pointer crossed the window, relative to the receiving window.
	//For Configure and Create requests, the x and y coordinates of the window relative to its parent window.
	PosX int
	PosY int

	GlobalPosX int
	GlobalPosY int

	//This reports the delta value of a MouseWheel event.
	//The delta value represents the rotation units the mouse wheel has been moved.
	//The sign of the value represents the direction the mouse wheel was scrolled.
	WheelDelta int

	//The keycode field from the event. Valid only for KeyPress and KeyRelease events.
	KeyCode int
	KeySym  string
	KeyText string
	KeyRune rune

	//The detail or user_data field from the event. The %d is replaced by a string identifying the detail.
	//For Enter, Leave, FocusIn, and FocusOut events, the string will be one of the following:
	//NotifyAncestor NotifyNonlinearVirtual NotifyDetailNone
	//NotifyPointer NotifyInferior NotifyPointerRoot NotifyNonlinear NotifyVirtual
	//For ConfigureRequest events, the string will be one of:
	//Above Opposite Below None BottomIf TopIf
	//For virtual events, the string will be whatever value is stored in the user_data field when the event was created (typically with event generate),
	//or the empty string if the field is NULL. Virtual events corresponding to key sequence presses (see event add for details) set the user_data to NULL.
	//For events other than these, the substituted string is undefined.
	UserData string

	//The focus field from the event (0 or 1). Valid only for Enter and Leave events.
	//1 if the receiving window is the focus window or a descendant of the focus window, 0 otherwise.
	Focus bool

	//The width/height field from the event. Valid for the Configure, ConfigureRequest, Create, ResizeRequest, and Expose events.
	//Indicates the new or requested width/height of the window.
	Width  int
	Height int

	//The mode field from the event. The substituted string is one of NotifyNormal, NotifyGrab, NotifyUngrab, or NotifyWhileGrabbed.
	//Valid only for Enter, FocusIn, FocusOut, and Leave events.
	Mode string

	//The override_redirect field from the event. Valid only for Map, Reparent, and Configure events.
	OverrideRedirect string

	//The place field from the event, substituted as one of the strings PlaceOnTop or PlaceOnBottom.
	//Valid only for Circulate and CirculateRequest events.
	Place string

	//The state field from the event. For ButtonPress, ButtonRelease, Enter, KeyPress, KeyRelease, Leave, and Motion events,
	//a decimal string is substituted. For Visibility, one of the strings VisibilityUnobscured, VisibilityPartiallyObscured, and VisibilityFullyObscured is substituted.
	//For Property events, substituted with either the string NewValue (indicating that the property has been created or modified) or Delete (indicating that the property has been removed).
	State string
}

func (e *Event) params() string {
	return "%T %E %W %t %b %x %y %D %k %K %A %d %f %w %h %m %o %p %s %X %Y"
}

func (e *Event) parser(args []string) {
	e.Type = e.toInt(args[0])
	e.Synthetic = e.toBool(args[1])
	e.Widget = FindWidget(args[2])
	e.Timestamp = e.toInt64(args[3])
	if e.Timestamp < 0 {
		e.Timestamp = 0
	}
	e.MouseButton = e.toInt(args[4])
	e.PosX = e.toInt(args[5])
	e.PosY = e.toInt(args[6])
	e.WheelDelta = e.toInt(args[7])
	e.KeyCode = e.toInt(args[8])
	e.KeySym = e.toString(args[9])
	e.KeyText = e.toString(args[10])
	if e.KeyText != "" {
		e.KeyRune, _ = utf8.DecodeRuneInString(e.KeyText)
	}
	e.UserData = e.toString(args[11])
	e.Focus = e.toBool(args[12])
	e.Width = e.toInt(args[13])
	e.Height = e.toInt(args[14])
	e.Mode = e.toString(args[15])
	e.OverrideRedirect = e.toString(args[16])
	e.Place = e.toString(args[17])
	e.State = e.toString(args[18])
	e.GlobalPosX = e.toInt(args[19])
	e.GlobalPosY = e.toInt(args[20])
}

func (e *Event) toInt(s string) int {
	v, _ := strconv.ParseInt(s, 10, 0)
	return int(v)
}

func (e *Event) toInt64(s string) int64 {
	v, _ := strconv.ParseInt(s, 10, 0)
	return v
}

func (e *Event) toBool(s string) bool {
	if s == "1" {
		return true
	}
	return false
}

func (e *Event) toString(s string) string {
	if s == "??" {
		return ""
	}
	return s
}

type KeyModifier int

const (
	KeyModifierNone KeyModifier = 1 << iota
	KeyModifierShift
	KeyModifierControl
	KeyModifierAlt
	KeyModifierMeta
	KeyModifierFn
)

func (k KeyModifier) String() string {
	var ar []string
	if k&KeyModifierShift == KeyModifierShift {
		ar = append(ar, "Shift")
	}
	if k&KeyModifierControl == KeyModifierControl {
		ar = append(ar, "Control")
	}
	if k&KeyModifierAlt == KeyModifierAlt {
		ar = append(ar, "Alt")
	}
	if k&KeyModifierMeta == KeyModifierMeta {
		ar = append(ar, "Meta")
	}
	return strings.Join(ar, " ")
}

type KeyEvent struct {
	*Event
	KeyModifier KeyModifier
}

func (e *KeyEvent) addModifier(sym string, name string, mod KeyModifier) {
	if strings.HasPrefix(sym, name) {
		e.KeyModifier |= mod
	}
}

func (e *KeyEvent) removeModifier(sym string, name string, mod KeyModifier) {
	if strings.HasPrefix(sym, name) {
		if e.KeyModifier&mod == mod {
			e.KeyModifier ^= mod
		}
	}
}

//TODO: almost check key modifier
func BindKeyEventEx(tag string, fnPress func(e *KeyEvent), fnRelease func(e *KeyEvent)) error {
	var ke KeyEvent
	var err error
	err = BindEvent(tag, "<KeyPress>", func(e *Event) {
		ke.addModifier(e.KeySym, "Shift_", KeyModifierShift)
		ke.addModifier(e.KeySym, "Control_", KeyModifierControl)
		ke.addModifier(e.KeySym, "Alt_", KeyModifierAlt)
		ke.addModifier(e.KeySym, "Meta_", KeyModifierMeta)
		ke.addModifier(e.KeySym, "Super_", KeyModifierFn)
		ke.Event = e
		if fnPress != nil {
			fnPress(&ke)
		}
	})
	if err != nil {
		return err
	}
	err = BindEvent(tag, "<KeyRelease>", func(e *Event) {
		ke.Event = e
		if fnRelease != nil {
			fnRelease(&ke)
		}
		ke.removeModifier(e.KeySym, "Shift_", KeyModifierShift)
		ke.removeModifier(e.KeySym, "Control_", KeyModifierControl)
		ke.removeModifier(e.KeySym, "Alt_", KeyModifierAlt)
		ke.removeModifier(e.KeySym, "Meta_", KeyModifierMeta)
		ke.removeModifier(e.KeySym, "Super_", KeyModifierFn)
	})
	return err
}

func bindEventHelper(tag string, event string, fnid string, ev *Event, fn func()) error {
	mainInterp.CreateAction(fnid, func(args []string) {
		ev.parser(args)
		fn()
	})
	return eval(fmt.Sprintf("bind %v %v {%v %v}", tag, event, fnid, ev.params()))
}

func addEventHelper(tag string, event string, fnid string, ev *Event, fn func()) error {
	mainInterp.CreateAction(fnid, func(args []string) {
		ev.parser(args)
		fn()
	})
	return eval(fmt.Sprintf("bind %v %v {+%v %v}", tag, event, fnid, ev.params()))
}

func IsEvent(event string) bool {
	return strings.HasPrefix(event, "<") && strings.HasSuffix(event, ">")
}

func IsVirtualEvent(event string) bool {
	return strings.HasPrefix(event, "<<") && strings.HasSuffix(event, ">>")
}

// add bind event
func BindEvent(tag string, event string, fn func(e *Event)) error {
	if tag == "" || !IsEvent(event) {
		return ErrInvalid
	}
	fnid := makeBindEventId()
	var ev Event
	return addEventHelper(tag, event, fnid, &ev, func() {
		fn(&ev)
	})
}

// clear tag event
func ClearBindEvent(tag string, event string) error {
	if tag == "" || !IsEvent(event) {
		return ErrInvalid
	}
	return eval(fmt.Sprintf("bind %v %v {}", tag, event))
}

func BindInfo(tag string) []string {
	if tag == "" {
		return nil
	}
	v, _ := evalAsStringList(fmt.Sprintf("bind %v", tag))
	return v
}

//Associates the virtual event virtual with the physical event sequence(s)
//given by the sequence arguments, so that the virtual event will trigger
//whenever any one of the sequences occurs. Virtual may be any string value
//and sequence may have any of the values allowed for the sequence argument
// to the bind command. If virtual is already defined, the new physical event
// sequences add to the existing sequences for the event.
func AddVirtualEventPhysicalEvent(virtual string, event string, events ...string) error {
	if !IsVirtualEvent(virtual) {
		return ErrInvalid
	}
	eventList := append([]string{event}, events...)
	return eval(fmt.Sprintf("event add %v %v", virtual, strings.Join(eventList, " ")))
}

//Deletes each of the sequences from those associated with the virtual event
// given by virtual. Virtual may be any string value and sequence may have
//any of the values allowed for the sequence argument to the bind command.
//Any sequences not currently associated with virtual are ignored.
//If no sequence argument is provided, all physical event sequences are removed
// for virtual, so that the virtual event will not trigger anymore.
func RemoveVirtualEventPhysicalEvent(virtual string, events ...string) error {
	if !IsVirtualEvent(virtual) {
		return ErrInvalid
	}
	return eval(fmt.Sprintf("event remove %v %v", virtual, strings.Join(events, " ")))
}

func VirtualEventInfo(virtual string) []string {
	if !IsVirtualEvent(virtual) {
		return nil
	}
	r, _ := evalAsStringList(fmt.Sprintf("event info %v", virtual))
	return r
}

//TODO: event attr
type EventAttr struct {
	key   string
	value string
}

func NativeEventAttr(key string, value string) *EventAttr {
	return &EventAttr{key, value}
}

func SendEvent(widget Widget, event string, attrs ...*EventAttr) error {
	if !IsValidWidget(widget) {
		return ErrInvalid
	}
	return sendEvent(widget.Id(), event, attrs...)
}

func SendEventToFocus(event string, attrs ...*EventAttr) error {
	return sendEvent("[focus]", event, attrs...)
}

func sendEvent(id string, event string, attrs ...*EventAttr) error {
	if !IsEvent(event) {
		return ErrInvalid
	}
	var list []string
	for _, attr := range attrs {
		if attr == nil {
			continue
		}
		list = append(list, fmt.Sprintf("-%v {%v}", attr.key, attr.value))
	}
	var script string
	script = fmt.Sprintf("event generate %v %v", id, event)
	if len(list) > 0 {
		script += " " + strings.Join(list, " ")
	}
	return eval(script)
}
