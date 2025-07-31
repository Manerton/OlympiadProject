<?php

namespace App\Builder;

use App\Models\Application;
use App\Models\Event;
use App\Models\User;

class ApplicationBuilder
{
    public function build($item){
        $model = new Application();
        $model->id = $item['id'];
        $model->user_id = $item['userId'];
        $model->code = $item['code'];
        $model->reason = $item['reason'];
        $model->event_id = $item['eventId'];
        $model->status = $item['status'];
        return $model;
    }
    public function buildEvent(Application $application, Event $event)
    {
        $application->eventAPI = $event;
    }
    public function buildUser(Application $application, User $user){
        $application->userAPI = $user;
    }
}
