<?php

namespace App\Builder;

use App\Models\Event;
use App\Models\EventJury;
use App\Models\User;

class EventJuryBuilder
{
    public function build($item){
        $model = new EventJury();
        $model->user_id = $item['user_id'];
        $model->event_id = $item['event_id'];
        return $model;
    }
    public function buildEvent(EventJury $eventJury, Event $event)
    {
        $eventJury->eventAPI = $event;
    }
    public function buildUser(EventJury $eventJury, User $user){
        $eventJury->userAPI = $user;
    }
}
