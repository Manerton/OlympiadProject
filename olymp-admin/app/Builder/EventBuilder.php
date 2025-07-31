<?php

namespace App\Builder;

use App\Models\Event;
use DateTime;

class EventBuilder
{
    public function build($item)
    {
        $event = new Event();
        $event->id = $item['id'];
        $event->name = $item['name'];
        $date = DateTime::createFromFormat(DateTime::ATOM,  $item['start_date']);
        $event->start_date  = $date->format('Y-m-d');
        $date = DateTime::createFromFormat(DateTime::ATOM, $item['end_date']);
        $event->end_date = $date->format('Y-m-d');
        $event->event_type = $item['event_type'];
        $event->class_number = $item['class_number'];
        $event->previous_event_id = $item['previous_event_id'];
        $event->subject = $item['subject'];
        $event->additional_info = $item['additional_info'];
        $event->events = $item['events'];
        return $event;
    }
}
