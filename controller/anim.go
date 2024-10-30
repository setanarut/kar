package controller

import (
	"kar/types"

	"github.com/setanarut/anim"
)

func UpdateAnimationStates(anim *anim.AnimationPlayer, opt *types.DrawOptions) {
	if isRunning {
		anim.SetStateFPS("walk_right", 25)
	} else {
		anim.SetStateFPS("walk_right", 15)
	}
	if isIdle && IsFacingLeft {
		anim.SetState("idle_left")
		opt.FlipX = false
	} else if isIdle && IsFacingRight {
		anim.SetState("idle_right")
		opt.FlipX = false
	} else {
		anim.SetState("idle_front")
		opt.FlipX = false
	}
	if !isOnFloor && !isIdle {
		anim.SetState("jump")
		if IsFacingLeft {
			opt.FlipX = true
		} else {
			opt.FlipX = false
		}
	}
	if isDigDown && !IsFacingLeft {
		anim.SetState("dig_down")
	}
	if isDigUp {
		anim.SetState("dig_right")
	}
	// dig right
	if isAttacking && IsFacingRight {
		anim.SetState("dig_right")
		opt.FlipX = false
	}
	// dig left
	if isAttacking && IsFacingLeft {
		anim.SetState("dig_right")
		opt.FlipX = true
	}
	if isWalkingRight {
		anim.SetState("walk_right")
		opt.FlipX = false
	}
	if isWalkingLeft {
		anim.SetState("walk_right")
		opt.FlipX = true
	}
}
