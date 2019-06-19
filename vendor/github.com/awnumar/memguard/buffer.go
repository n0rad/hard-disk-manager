package memguard

import (
	"runtime"
	"unsafe"

	"github.com/awnumar/memguard/core"
)

/*
LockedBuffer is a structure that holds raw sensitive data.

The number of LockedBuffers that you are able to create is limited by how much memory your system's kernel allows each process to mlock/VirtualLock. Therefore you should call Destroy on LockedBuffers that you no longer need or defer a Destroy call after creating a new LockedBuffer.
*/
type LockedBuffer struct {
	*core.Buffer
	*drop
}

/*
This is a value that is monitored by a finalizer so that we can clean up LockedBuffers that have gone out of scope.
*/
type drop [16]byte

/*
NewBuffer creates a mutable data container of the specified size.

The size must be strictly positive or the function will panic.
*/
func NewBuffer(size int) *LockedBuffer {
	// Construct a Buffer of the specified size.
	buf, err := core.NewBuffer(size)
	if err != nil {
		core.Panic(err)
	}

	// Initialise a LockedBuffer object around it.
	b := &LockedBuffer{buf, new(drop)}

	// Use a finalizer to destroy the Buffer if it falls out of scope.
	runtime.SetFinalizer(b.drop, func(_ *drop) {
		go buf.Destroy()
	})

	// Return the created buffer to the caller.
	return b
}

/*
NewBufferFromBytes constructs an immutable buffer from a byte slice.

The length of the buffer must be strictly positive or the function will panic. The source buffer is wiped after the value has been copied over to the created container.
*/
func NewBufferFromBytes(src []byte) *LockedBuffer {
	// Construct a buffer of the correct size.
	b := NewBuffer(len(src))

	// Move the data over.
	core.Move(b.Bytes(), src)

	// Make the buffer immutable.
	b.Freeze()

	// Return the created Buffer object.
	return b
}

/*
NewBufferRandom constructs an immutable buffer filled with cryptographically-secure random bytes.

The size must be strictly positive or the function will panic.
*/
func NewBufferRandom(size int) *LockedBuffer {
	// Construct a buffer of the specified size.
	b := NewBuffer(size)

	// Fill the buffer with random bytes.
	ScrambleBytes(b.Bytes())

	// Make the buffer immutable.
	b.Freeze()

	// Return the created Buffer object.
	return b
}

// Freeze makes a LockedBuffer's memory immutable. The call can be reversed with Melt.
func (b *LockedBuffer) Freeze() {
	b.Buffer.Freeze()
}

// Melt makes a LockedBuffer's memory mutable. The call can be reversed with Freeze.
func (b *LockedBuffer) Melt() {
	b.Buffer.Melt()
}

/*
Seal takes a LockedBuffer object and returns its contents encrypted inside a sealed Enclave object. The LockedBuffer is subsequently destroyed and its contents wiped.

If Seal is called on a destroyed buffer, a nil enclave is returned.
*/
func (b *LockedBuffer) Seal() *Enclave {
	e, err := core.Seal(b.Buffer)
	if err != nil {
		if err == core.ErrBufferExpired {
			return nil
		}
		core.Panic(err)
	}
	return &Enclave{e}
}

/*
Copy performs a time-constant copy into a LockedBuffer. Move is preferred if the source is not also a LockedBuffer or if the source is no longer needed.
*/
func (b *LockedBuffer) Copy(src []byte) {
	b.CopyAt(0, src)
}

/*
CopyAt performs a time-constant copy into a LockedBuffer at an offset. Move is preferred if the source is not also a LockedBuffer or if the source is no longer needed.
*/
func (b *LockedBuffer) CopyAt(offset int, src []byte) {
	if !b.IsAlive() {
		return
	}

	b.Lock()
	defer b.Unlock()

	core.Copy(b.Bytes()[offset:], src)
}

/*
Move performs a time-constant move into a LockedBuffer. The source is wiped after the bytes are copied.
*/
func (b *LockedBuffer) Move(src []byte) {
	b.MoveAt(0, src)
}

/*
MoveAt performs a time-constant move into a LockedBuffer at an offset. The source is wiped after the bytes are copied.
*/
func (b *LockedBuffer) MoveAt(offset int, src []byte) {
	if !b.IsAlive() {
		return
	}

	b.Lock()
	defer b.Unlock()

	core.Move(b.Bytes()[offset:], src)
}

/*
Scramble attempts to overwrite the data with cryptographically-secure random bytes.
*/
func (b *LockedBuffer) Scramble() {
	if !b.IsAlive() {
		return
	}

	b.Lock()
	defer b.Unlock()

	core.Scramble(b.Bytes())
}

/*
Wipe attempts to overwrite the data with zeros.
*/
func (b *LockedBuffer) Wipe() {
	if !b.IsAlive() {
		return
	}

	b.Lock()
	defer b.Unlock()

	core.Wipe(b.Bytes())
}

/*
Size gives you the length of a given LockedBuffer's data segment. A destroyed LockedBuffer will have a size of zero.
*/
func (b *LockedBuffer) Size() int {
	return len(b.Bytes())
}

/*
Destroy wipes and frees the underlying memory of a LockedBuffer. The LockedBuffer will not be accessible or usable after this calls is made.
*/
func (b *LockedBuffer) Destroy() {
	b.Buffer.Destroy()
}

/*
IsAlive returns a boolean value indicating if a LockedBuffer is alive, i.e. that it has not been destroyed.
*/
func (b *LockedBuffer) IsAlive() bool {
	return core.GetBufferState(b.Buffer).IsAlive
}

/*
IsMutable returns a boolean value indicating if a LockedBuffer is mutable.
*/
func (b *LockedBuffer) IsMutable() bool {
	return core.GetBufferState(b.Buffer).IsMutable
}

/*
EqualTo performs a time-constant comparison on the contents of a LockedBuffer with a given buffer. A destroyed LockedBuffer will always return false.
*/
func (b *LockedBuffer) EqualTo(buf []byte) bool {
	b.RLock()
	defer b.RUnlock()

	return core.Equal(b.Bytes(), buf)
}

/*
	Functions for representing the memory region as various data types.
*/

/*
Bytes returns a byte slice referencing the protected region of memory.
*/
func (b *LockedBuffer) Bytes() []byte {
	return b.Buffer.Data()
}

/*
Uint16 returns a slice pointing to the protected region of memory with the data represented as a sequence of unsigned 16 bit integers. Its length will be half that of the byte slice, excluding any remaining part that doesn't form a complete uint16 value.

If called on a destroyed LockedBuffer, a nil slice will be returned.
*/
func (b *LockedBuffer) Uint16() []uint16 {
	b.RLock()
	defer b.RUnlock()

	// Check if still alive.
	if !core.GetBufferState(b.Buffer).IsAlive {
		return nil
	}

	// Compute size of new slice representation.
	size := b.Size() / 2
	if size < 1 {
		return nil
	}

	// Construct the new slice representation.
	var sl = struct {
		addr uintptr
		len  int
		cap  int
	}{uintptr(unsafe.Pointer(&b.Bytes()[0])), size, size}

	// Cast the representation to the correct type and return it.
	return *(*[]uint16)(unsafe.Pointer(&sl))
}

/*
Uint32 returns a slice pointing to the protected region of memory with the data represented as a sequence of unsigned 32 bit integers. Its length will be one quarter that of the byte slice, excluding any remaining part that doesn't form a complete uint32 value.

If called on a destroyed LockedBuffer, a nil slice will be returned.
*/
func (b *LockedBuffer) Uint32() []uint32 {
	b.RLock()
	defer b.RUnlock()

	// Check if still alive.
	if !core.GetBufferState(b.Buffer).IsAlive {
		return nil
	}

	// Compute size of new slice representation.
	size := b.Size() / 4
	if size < 1 {
		return nil
	}

	// Construct the new slice representation.
	var sl = struct {
		addr uintptr
		len  int
		cap  int
	}{uintptr(unsafe.Pointer(&b.Bytes()[0])), size, size}

	// Cast the representation to the correct type and return it.
	return *(*[]uint32)(unsafe.Pointer(&sl))
}

/*
Uint64 returns a slice pointing to the protected region of memory with the data represented as a sequence of unsigned 64 bit integers. Its length will be one eighth that of the byte slice, excluding any remaining part that doesn't form a complete uint64 value.

If called on a destroyed LockedBuffer, a nil slice will be returned.
*/
func (b *LockedBuffer) Uint64() []uint64 {
	b.RLock()
	defer b.RUnlock()

	// Check if still alive.
	if !core.GetBufferState(b.Buffer).IsAlive {
		return nil
	}

	// Compute size of new slice representation.
	size := b.Size() / 8
	if size < 1 {
		return nil
	}

	// Construct the new slice representation.
	var sl = struct {
		addr uintptr
		len  int
		cap  int
	}{uintptr(unsafe.Pointer(&b.Bytes()[0])), size, size}

	// Cast the representation to the correct type and return it.
	return *(*[]uint64)(unsafe.Pointer(&sl))
}

/*
Int8 returns a slice pointing to the protected region of memory with the data represented as a sequence of signed 8 bit integers. If called on a destroyed LockedBuffer, a nil slice will be returned.
*/
func (b *LockedBuffer) Int8() []int8 {
	b.RLock()
	defer b.RUnlock()

	// Check if still alive.
	if !core.GetBufferState(b.Buffer).IsAlive {
		return nil
	}

	// Construct the new slice representation.
	var sl = struct {
		addr uintptr
		len  int
		cap  int
	}{uintptr(unsafe.Pointer(&b.Bytes()[0])), b.Size(), b.Size()}

	// Cast the representation to the correct type and return it.
	return *(*[]int8)(unsafe.Pointer(&sl))
}

/*
Int16 returns a slice pointing to the protected region of memory with the data represented as a sequence of signed 16 bit integers. Its length will be half that of the byte slice, excluding any remaining part that doesn't form a complete int16 value.

If called on a destroyed LockedBuffer, a nil slice will be returned.
*/
func (b *LockedBuffer) Int16() []int16 {
	b.RLock()
	defer b.RUnlock()

	// Check if still alive.
	if !core.GetBufferState(b.Buffer).IsAlive {
		return nil
	}

	// Compute size of new slice representation.
	size := b.Size() / 2
	if size < 1 {
		return nil
	}

	// Construct the new slice representation.
	var sl = struct {
		addr uintptr
		len  int
		cap  int
	}{uintptr(unsafe.Pointer(&b.Bytes()[0])), size, size}

	// Cast the representation to the correct type and return it.
	return *(*[]int16)(unsafe.Pointer(&sl))
}

/*
Int32 returns a slice pointing to the protected region of memory with the data represented as a sequence of signed 32 bit integers. Its length will be one quarter that of the byte slice, excluding any remaining part that doesn't form a complete int32 value.

If called on a destroyed LockedBuffer, a nil slice will be returned.
*/
func (b *LockedBuffer) Int32() []int32 {
	b.RLock()
	defer b.RUnlock()

	// Check if still alive.
	if !core.GetBufferState(b.Buffer).IsAlive {
		return nil
	}

	// Compute size of new slice representation.
	size := b.Size() / 4
	if size < 1 {
		return nil
	}

	// Construct the new slice representation.
	var sl = struct {
		addr uintptr
		len  int
		cap  int
	}{uintptr(unsafe.Pointer(&b.Bytes()[0])), size, size}

	// Cast the representation to the correct type and return it.
	return *(*[]int32)(unsafe.Pointer(&sl))
}

/*
Int64 returns a slice pointing to the protected region of memory with the data represented as a sequence of signed 64 bit integers. Its length will be one eighth that of the byte slice, excluding any remaining part that doesn't form a complete int64 value.

If called on a destroyed LockedBuffer, a nil slice will be returned.
*/
func (b *LockedBuffer) Int64() []int64 {
	b.RLock()
	defer b.RUnlock()

	// Check if still alive.
	if !core.GetBufferState(b.Buffer).IsAlive {
		return nil
	}

	// Compute size of new slice representation.
	size := b.Size() / 8
	if size < 1 {
		return nil
	}

	// Construct the new slice representation.
	var sl = struct {
		addr uintptr
		len  int
		cap  int
	}{uintptr(unsafe.Pointer(&b.Bytes()[0])), size, size}

	// Cast the representation to the correct type and return it.
	return *(*[]int64)(unsafe.Pointer(&sl))
}

/*
ByteArray8 returns a pointer to some 8 byte array. Care must be taken not to dereference the pointer and instead pass it around as-is.

The length of the buffer must be at least 8 bytes in size and the LockedBuffer should not be destroyed. In either of these cases a nil value is returned.
*/
func (b *LockedBuffer) ByteArray8() *[8]byte {
	b.RLock()
	defer b.RUnlock()

	// Check if still alive.
	if !core.GetBufferState(b.Buffer).IsAlive {
		return nil
	}

	// Check if the length is large enough.
	if len(b.Bytes()) < 8 {
		return nil
	}

	// Cast the representation to the correct type.
	return (*[8]byte)(unsafe.Pointer(&b.Bytes()[0]))
}

/*
ByteArray16 returns a pointer to some 16 byte array. Care must be taken not to dereference the pointer and instead pass it around as-is.

The length of the buffer must be at least 16 bytes in size and the LockedBuffer should not be destroyed. In either of these cases a nil value is returned.
*/
func (b *LockedBuffer) ByteArray16() *[16]byte {
	b.RLock()
	defer b.RUnlock()

	// Check if still alive.
	if !core.GetBufferState(b.Buffer).IsAlive {
		return nil
	}

	// Check if the length is large enough.
	if len(b.Bytes()) < 16 {
		return nil
	}

	// Cast the representation to the correct type.
	return (*[16]byte)(unsafe.Pointer(&b.Bytes()[0]))
}

/*
ByteArray32 returns a pointer to some 32 byte array. Care must be taken not to dereference the pointer and instead pass it around as-is.

The length of the buffer must be at least 32 bytes in size and the LockedBuffer should not be destroyed. In either of these cases a nil value is returned.
*/
func (b *LockedBuffer) ByteArray32() *[32]byte {
	b.RLock()
	defer b.RUnlock()

	// Check if still alive.
	if !core.GetBufferState(b.Buffer).IsAlive {
		return nil
	}

	// Check if the length is large enough.
	if len(b.Bytes()) < 32 {
		return nil
	}

	// Cast the representation to the correct type.
	return (*[32]byte)(unsafe.Pointer(&b.Bytes()[0]))
}

/*
ByteArray64 returns a pointer to some 64 byte array. Care must be taken not to dereference the pointer and instead pass it around as-is.

The length of the buffer must be at least 64 bytes in size and the LockedBuffer should not be destroyed. In either of these cases a nil value is returned.
*/
func (b *LockedBuffer) ByteArray64() *[64]byte {
	b.RLock()
	defer b.RUnlock()

	// Check if still alive.
	if !core.GetBufferState(b.Buffer).IsAlive {
		return nil
	}

	// Check if the length is large enough.
	if len(b.Bytes()) < 64 {
		return nil
	}

	// Cast the representation to the correct type.
	return (*[64]byte)(unsafe.Pointer(&b.Bytes()[0]))
}
