package mqtt

import (
	"strings"
	"testing"
	"time"
)

func TestSubscribeNotConnected(t *testing.T) {
	ops := NewClientOptions().SetClientID("TestSubscribeNotConnected")
	c := NewClient(ops)

	token := c.Subscribe("test/topic", 0, nil)
	if token.Wait() && token.Error() == nil {
		t.Fatal("Subscribe should fail when not connected")
	}
	if !strings.Contains(token.Error().Error(), "订阅失败") {
		t.Fatalf("错误信息应该包含'订阅失败'，实际错误信息: %v", token.Error())
	}
}

func TestSubscribeNotConnectedNoResume(t *testing.T) {
	ops := NewClientOptions().SetClientID("TestSubscribeNotConnectedNoResume")
	ops.SetResumeSubs(false)
	c := NewClient(ops)

	// 模拟连接状态为未连接
	cli := c.(*client)
	cli.status.forceConnectionStatus(disconnected)

	token := c.Subscribe("test/topic", 0, nil)
	if token.Wait() && token.Error() == nil {
		t.Fatal("Subscribe should fail when not connected and ResumeSubs is false")
	}
	if !strings.Contains(token.Error().Error(), "订阅失败") {
		t.Fatalf("错误信息应该包含'订阅失败'，实际错误信息: %v", token.Error())
	}
}

func TestSubscribeReconnectingCleanSession(t *testing.T) {
	ops := NewClientOptions().SetClientID("TestSubscribeReconnectingCleanSession")
	ops.SetCleanSession(true)
	c := NewClient(ops)

	// 模拟重连状态
	cli := c.(*client)
	cli.status.forceConnectionStatus(reconnecting)

	token := c.Subscribe("test/topic", 0, nil)
	if token.Wait() && token.Error() == nil {
		t.Fatal("Subscribe should fail when reconnecting and CleanSession is true")
	}
	if !strings.Contains(token.Error().Error(), "订阅失败") {
		t.Fatalf("错误信息应该包含'订阅失败'，实际错误信息: %v", token.Error())
	}
}

func TestSubscribeInvalidTopic(t *testing.T) {
	ops := NewClientOptions().SetClientID("TestSubscribeInvalidTopic")
	c := NewClient(ops)

	// 模拟连接状态为已连接
	cli := c.(*client)
	cli.status.forceConnectionStatus(connected)

	token := c.Subscribe("", 0, nil) // 空topic
	if token.Wait() && token.Error() == nil {
		t.Fatal("Subscribe should fail with invalid topic")
	}
	if !strings.Contains(token.Error().Error(), "订阅失败") {
		t.Fatalf("错误信息应该包含'订阅失败'，实际错误信息: %v", token.Error())
	}
}

func TestSubscribeInvalidQos(t *testing.T) {
	ops := NewClientOptions().SetClientID("TestSubscribeInvalidQos")
	c := NewClient(ops)

	// 模拟连接状态为已连接
	cli := c.(*client)
	cli.status.forceConnectionStatus(connected)

	token := c.Subscribe("test/topic", 3, nil) // 无效的QoS
	if token.Wait() && token.Error() == nil {
		t.Fatal("Subscribe should fail with invalid QoS")
	}
	if !strings.Contains(token.Error().Error(), "订阅失败") {
		t.Fatalf("错误信息应该包含'订阅失败'，实际错误信息: %v", token.Error())
	}
}

func TestSubscribeNoMessageID(t *testing.T) {
	ops := NewClientOptions().SetClientID("TestSubscribeNoMessageID")
	c := NewClient(ops)

	// 模拟连接状态为已连接
	cli := c.(*client)
	cli.status.forceConnectionStatus(connected)

	// 模拟消息ID不足
	cli.messageIds = messageIds{index: make(map[uint16]tokenCompletor)}
	for i := uint16(1); i <= 65535; i++ {
		cli.messageIds.index[i] = nil
	}

	token := c.Subscribe("test/topic", 0, nil)
	if token.Wait() && token.Error() == nil {
		t.Fatal("Subscribe should fail when no message IDs available")
	}
	if !strings.Contains(token.Error().Error(), "订阅失败") {
		t.Fatalf("错误信息应该包含'invalid QoS'，实际错误信息: %v", token.Error())
	}
}

func TestSubscribeTimeout(t *testing.T) {
	ops := NewClientOptions().SetClientID("TestSubscribeTimeout")
	ops.SetWriteTimeout(1 * time.Millisecond) // 设置很短的超时时间
	c := NewClient(ops)

	// 模拟连接状态为已连接
	cli := c.(*client)
	cli.status.forceConnectionStatus(connected)

	// 模拟oboundP通道已满
	cli.oboundP = make(chan *PacketAndToken, 0)

	token := c.Subscribe("test/topic", 0, nil)
	if token.Wait() && token.Error() == nil {
		t.Fatal("Subscribe should fail with timeout")
	}
	if !strings.Contains(token.Error().Error(), "订阅失败") {
		t.Fatalf("错误信息应该包含'invalid QoS'，实际错误信息: %v", token.Error())
	}
}
