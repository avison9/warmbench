// Package driver defines the engine-agnostic protocol between the daemon and
// a slot (JSON lines over a Unix socket or named pipe) and the Unity driver
// that speaks it to the com.warmbench.agent editor package.
//
// The scheduler, store and MCP layer must not import Unity-specific types;
// other engines are added as further drivers.
package driver
