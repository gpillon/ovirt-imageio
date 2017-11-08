"""
Handling the proxy /tickets/ resource.
"""

from six.moves import http_client
from webob import exc
from ovirt_imageio_common import web

from . import auth2


class RequestHandler(object):
    """
    Request handler for the /tickets/ resource.
    """

    def __init__(self, config, request):
        """
        Arguments:
            config (config object): proxy configuration
            request (webob.Request): underlying http request
        """
        self.config = config
        self.request = request

    def put(self, ticket_id=None):
        """
        Verify and add a signed_ticket, allowing transfer to /images/ticket_id.
        """
        try:
            auth2.add_signed_ticket(self.request.body)
        except auth2.Error as e:
            raise exc.HTTPForbidden("Error verifying signed ticket: %s" % e)
        return web.response()

    def delete(self, ticket_id=None):
        """
        Delete ticket by a specified ticket_id
        """
        if not ticket_id:
            raise exc.HTTPBadRequest("Missing ticket ID")
        try:
            auth2.delete_ticket(ticket_id)
        except auth2.NoSuchTicket as e:
            raise exc.HTTPNotFound("Ticket not found: %s" % e)
        return web.response(http_client.NO_CONTENT)